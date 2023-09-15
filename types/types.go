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
