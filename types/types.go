package types

import (
	"time"

	"github.com/google/uuid"

	"github.com/adhupraba/discord-server/internal/discord/public/model"
)

type Json map[string]any

type PaginationOpts struct {
	Limit  int64
	Offset int64
}

type ProfileAndServer struct {
	model.Profiles
	Server *model.Servers `json:"server"`
}

type MemberWithProfile struct {
	model.Members
	Profile model.Profiles `json:"profile"`
}

type ServerWithChannelsAndMembers struct {
	model.Servers
	Channels []model.Channels    `json:"channels"`
	Members  []MemberWithProfile `json:"members"`
}

type ConversationWithMemberAndProfile struct {
	model.Conversations
	MemberOne struct {
		model.Members `alias:"member_one.*"`
		Profile       model.Profiles `alias:"profile_one.*" json:"profile"`
	} `json:"memberOne"`
	MemberTwo struct {
		model.Members `alias:"member_two.*"`
		Profile       model.Profiles `alias:"profile_two.*" json:"profile"`
	} `json:"memberTwo"`
}

type ServerWithMembers struct {
	model.Servers
	Members []model.Members `json:"members"`
}

type DbMessageWithMember struct {
	model.Messages
	Member MemberWithProfile `json:"member"`
}

type WsIncomingMessageBody struct {
	Content string `json:"content" validate:"required,min=1"`
	FileUrl string `json:"fileUrl" validate:"url"`
}

type WsMessageEvent string

const (
	WsMessageEventBROADCAST   WsMessageEvent = "BROADCAST"
	WsMessageEventSENDMESSAGE WsMessageEvent = "SENDMESSAGE"
)

type WsIncomingMessage struct {
	Event   WsMessageEvent         `json:"event" validate:"required,min=1"`
	UserID  string                 `json:"userId" validate:"required,min=1"`
	Message *WsIncomingMessageBody `json:"message"`
}

type WsMessage struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	FileUrl   *string   `json:"fileUrl"`
	MemberID  uuid.UUID `json:"memberId"`
	RoomId    uuid.UUID `json:"roomId"` // kept as generic. can be either channel id or conversation id
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WsMessageContent struct {
	WsMessage
	Member MemberWithProfile `json:"member"`
}

type WsOutgoingMessage struct {
	Event   WsMessageEvent    `json:"event"`
	UserID  string            `json:"userId"` // user id -> (profile table id)
	Message *WsMessageContent `json:"message"`
}
