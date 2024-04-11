package types

import "github.com/adhupraba/discord-server/internal/discord/public/model"

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

type SendMessageBody struct {
	Content string `json:"content" validate:"required,min=1"`
	FileUrl string `json:"fileUrl" validate:"url"`
}

type MessageWithMember struct {
	model.Messages
	Member MemberWithProfile `json:"member"`
}
