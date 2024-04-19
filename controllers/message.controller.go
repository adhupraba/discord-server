package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/helpers"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type MessageController struct{}

type GetMessagesRes struct {
	NextCursor *string                  `json:"nextCursor"`
	Messages   []types.WsMessageContent `json:"messages"`
}

func (mc *MessageController) GetMessages(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	cursor := r.URL.Query().Get("cursor")
	var lastMsgId *uuid.UUID = nil
	var lastMsgDate *time.Time = nil

	// cursor is base64 version of `{{id}}&{{date}}`
	// eg: decoded cursor -> "6d35b78d-a9ef-44dc-bd6e-d0b60c39a0e0&2024-04-18T18:38:53.739Z"
	// encoded cursor -> "NmQzNWI3OGQtYTllZi00NGRjLWJkNmUtZDBiNjBjMzlhMGUwJjIwMjQtMDQtMThUMTg6Mzg6NTMuNzM5Wg=="
	if cursor != "" {
		decoded, err := helpers.Base64Decode(cursor)

		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid cursor")
			return
		}

		split := strings.Split(decoded, "&")

		if len(split) != 2 {
			utils.RespondWithError(w, http.StatusBadRequest, "Malformed cursor")
			return
		}

		uuidParsed, err := uuid.Parse(split[0])

		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid message id")
			return
		} else {
			lastMsgId = &uuidParsed
		}

		t, err := time.Parse(time.RFC3339, split[1])

		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid date")
			return
		}

		lastMsgDate = &t
	}

	channelIdQuery := r.URL.Query().Get("channelId")

	if channelIdQuery == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Channel ID missing")
		return
	}

	channelId, err := uuid.Parse(channelIdQuery)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid channel id")
		return
	}

	messages, nextCursor, err := lib.DB.GetMessages(r.Context(), queries.GetMessagesParams{
		ChannelId:       channelId,
		LastMessageId:   lastMsgId,
		LastMessageDate: lastMsgDate,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := GetMessagesRes{
		NextCursor: nextCursor,
		Messages:   messages,
	}

	utils.RespondWithJson(w, http.StatusOK, res)
}

func (mc *MessageController) SendMessage(w http.ResponseWriter, r *http.Request, user model.Profiles) {
	var body types.WsIncomingMessageBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data")
		return
	}

	channel, member, errCode, err := validateServer_Channel_Member(r.Context(), user.ID, r.URL.Query().Get("serverId"), r.URL.Query().Get("channelId"))

	if err != nil {
		utils.RespondWithError(w, errCode, err.Error())
		return
	}

	wsMessage, err := lib.BroadcastMessage(member.ID.String(), channel.ID.String(), types.WsRoomTypeCHANNEL, body)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJson(w, http.StatusOK, wsMessage)
}

func validateServer_Channel_Member(ctx context.Context, userId uuid.UUID, server_id string, channel_id string) (channel *model.Channels, member *model.Members, errCode int, e error) {
	serverId, err := uuid.Parse(server_id)

	if err != nil {
		return nil, nil, http.StatusBadRequest, errors.New("Invalid server id")
	}

	channeId, err := uuid.Parse(channel_id)

	if err != nil {
		return nil, nil, http.StatusBadRequest, errors.New("Invalid channel id")
	}

	server, err := lib.DB.GetServerAndMembersOfUser(ctx, queries.GetServerAndMembersOfUserParam{
		ServerId:  serverId,
		ProfileId: userId,
	})

	if err == qrm.ErrNoRows {
		return nil, nil, http.StatusNotFound, errors.New("Server not found")
	}

	if err != nil {
		log.Println("get server and members error:", err)
		return nil, nil, http.StatusInternalServerError, errors.New("Error getting server details")
	}

	chann, err := lib.DB.GetServerChannel(ctx, queries.GetServerChannelParams{
		ServerId:  &server.ID,
		ChannelId: channeId,
	})

	if err == qrm.ErrNoRows {
		return nil, nil, http.StatusNotFound, errors.New("Channel not found")
	}

	if err != nil {
		log.Println("get channel error:", err)
		return nil, nil, http.StatusInternalServerError, errors.New("Error getting channel details")
	}

	for _, sMem := range server.Members {
		if sMem.ProfileID.String() == userId.String() {
			member = &sMem
			break
		}
	}

	if member == nil {
		return nil, nil, http.StatusNotFound, errors.New("Member not found")
	}

	return &chann, member, 0, nil
}
