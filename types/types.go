package types

import "github.com/adhupraba/discord-server/internal/discord/public/model"

type Json map[string]any

type ProfileAndServers struct {
	model.Profiles
	Servers []model.Servers `json:"servers"`
}
