package controllers

import (
	"net/http"
	"strings"
	"time"

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
