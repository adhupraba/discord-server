package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type WsController struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wc *WsController) Connect(w http.ResponseWriter, r *http.Request, user model.Profiles) {
	fmt.Println("connect endpoint reacher for user =>", user.Email)

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error establishing websocket channel connection")
		return
	}

	fmt.Println("websocket conn established for user =>", user.Email)

	mem := &lib.WsClient{
		Conn:    conn,
		ID:      user.ID.String(),
		Message: make(chan *types.WsOutgoingMessage),
	}

	lib.WsHub.Register <- mem

	go mem.WriteMessage()

	mem.ReadMessage()
}

func (wc *WsController) SendMessage(w http.ResponseWriter, r *http.Request, user model.Profiles) {
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
