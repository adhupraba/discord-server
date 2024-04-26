package controllers

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/helpers"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/types"
	"github.com/adhupraba/discord-server/utils"
)

type DirectMessageController struct{}

type GetDirectMessagesRes struct {
	NextCursor *string                  `json:"nextCursor"`
	Messages   []types.WsMessageContent `json:"messages"`
}

func (mc *DirectMessageController) GetMessages(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	cursor := r.URL.Query().Get("cursor")

	lastMsgId, lastMsgDate, err := helpers.ValidateCursor(cursor)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	conversationId := r.URL.Query().Get("conversationId")

	if conversationId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Channel ID missing")
		return
	}

	conversationUUID, err := uuid.Parse(conversationId)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid conversation id")
		return
	}

	messages, nextCursor, err := lib.DB.GetDirectMessages(r.Context(), queries.GetDirectMessagesParams{
		ConversationID:  conversationUUID,
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
