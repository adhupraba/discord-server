package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
	"github.com/adhupraba/discord-server/internal/queries"
	"github.com/adhupraba/discord-server/lib"
	"github.com/adhupraba/discord-server/utils"
)

type ConversationController struct{}

func (cc *ConversationController) GetConversationByMembers(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	m1Id := r.URL.Query().Get("memberOne")
	m2Id := r.URL.Query().Get("memberTwo")

	memberOneId, err := uuid.Parse(m1Id)
	memberTwoId, err := uuid.Parse(m2Id)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid member id")
		return
	}

	conversation, err := lib.DB.GetConversationByMembers(r.Context(), queries.GetConversationByMembersParams{
		MemberOneId: memberOneId,
		MemberTwoId: memberTwoId,
	})

	if err == qrm.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Conversation not found")
		return
	}

	if err != nil {
		fmt.Println("Error when fetching conversation", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error when fetching conversation")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, conversation)
}

func (cc *ConversationController) CreateNewConversation(w http.ResponseWriter, r *http.Request, profile model.Profiles) {
	type createNewConversationBody struct {
		MemberOneId string `json:"memberOneId" validate:"required,uuid4"`
		MemberTwoId string `json:"memberTwoId" validate:"required,uuid4"`
	}

	var body createNewConversationBody
	err := utils.BodyParser(r.Body, &body)

	if err != nil {
		utils.RespondWithError(w, http.StatusUnprocessableEntity, "Invalid data")
		return
	}

	memberOneId, err := uuid.Parse(body.MemberOneId)
	memberTwoId, err := uuid.Parse(body.MemberTwoId)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid member id")
		return
	}

	_, err = lib.DB.CreateConversation(r.Context(), model.Conversations{
		ID:          uuid.New(),
		MemberOneID: memberOneId,
		MemberTwoID: memberTwoId,
	})

	if err != nil {
		log.Println("error creating conversation =>", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating conversation")
		return
	}

	conversation, err := lib.DB.GetConversationByMembers(r.Context(), queries.GetConversationByMembersParams{
		MemberOneId: memberOneId,
		MemberTwoId: memberTwoId,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching conversation data")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, conversation)
}
