package vlessproto

import "github.com/gofrs/uuid/v5"

type User struct {
	Desc string    `json:"desc"`
	UUID uuid.UUID `json:"uuid"`
	Flow string    `json:"flow,omitempty"`
}

type UserStore interface {
	GetVless(id uuid.UUID) (User, bool, error)
}
