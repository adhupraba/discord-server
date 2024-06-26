// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ChannelType string

const (
	ChannelTypeTEXT  ChannelType = "TEXT"
	ChannelTypeAUDIO ChannelType = "AUDIO"
	ChannelTypeVIDEO ChannelType = "VIDEO"
)

func (e *ChannelType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ChannelType(s)
	case string:
		*e = ChannelType(s)
	default:
		return fmt.Errorf("unsupported scan type for ChannelType: %T", src)
	}
	return nil
}

type NullChannelType struct {
	ChannelType ChannelType `json:"channelType"`
	Valid       bool        `json:"valid"` // Valid is true if ChannelType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullChannelType) Scan(value interface{}) error {
	if value == nil {
		ns.ChannelType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ChannelType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullChannelType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ChannelType), nil
}

type MemberRole string

const (
	MemberRoleADMIN     MemberRole = "ADMIN"
	MemberRoleMODERATOR MemberRole = "MODERATOR"
	MemberRoleGUEST     MemberRole = "GUEST"
)

func (e *MemberRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MemberRole(s)
	case string:
		*e = MemberRole(s)
	default:
		return fmt.Errorf("unsupported scan type for MemberRole: %T", src)
	}
	return nil
}

type NullMemberRole struct {
	MemberRole MemberRole `json:"memberRole"`
	Valid      bool       `json:"valid"` // Valid is true if MemberRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMemberRole) Scan(value interface{}) error {
	if value == nil {
		ns.MemberRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MemberRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMemberRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MemberRole), nil
}

type Channels struct {
	ID        uuid.UUID   `json:"id" sql:"primary_key"`
	Name      string      `json:"name"`
	Type      ChannelType `json:"type"`
	ProfileID uuid.UUID   `json:"profileId"`
	ServerID  uuid.UUID   `json:"serverId"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type Conversations struct {
	ID          uuid.UUID `json:"id" sql:"primary_key"`
	MemberOneID uuid.UUID `json:"memberOneId"`
	MemberTwoID uuid.UUID `json:"memberTwoId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DirectMessages struct {
	ID             uuid.UUID `json:"id" sql:"primary_key"`
	Content        string    `json:"content"`
	FileURL        *string   `json:"fileUrl"`
	MemberID       uuid.UUID `json:"memberId"`
	ConversationID uuid.UUID `json:"conversationId"`
	Deleted        bool      `json:"deleted"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Members struct {
	ID        uuid.UUID  `json:"id" sql:"primary_key"`
	Role      MemberRole `json:"role"`
	ProfileID uuid.UUID  `json:"profileId"`
	ServerID  uuid.UUID  `json:"serverId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type Messages struct {
	ID        uuid.UUID `json:"id" sql:"primary_key"`
	Content   string    `json:"content"`
	FileURL   *string   `json:"fileUrl"`
	MemberID  uuid.UUID `json:"memberId"`
	ChannelID uuid.UUID `json:"channelId"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Profiles struct {
	ID        uuid.UUID `json:"id" sql:"primary_key"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"imageUrl"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Servers struct {
	ID         uuid.UUID `json:"id" sql:"primary_key"`
	Name       string    `json:"name"`
	ImageURL   string    `json:"imageUrl"`
	InviteCode uuid.UUID `json:"inviteCode"`
	ProfileID  uuid.UUID `json:"profileId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
