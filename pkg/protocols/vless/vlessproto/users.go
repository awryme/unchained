package vlessproto

import "github.com/gofrs/uuid/v5"

type User struct {
	Name string    `json:"name"`
	UUID uuid.UUID `json:"uuid"`
	Flow string    `json:"flow,omitempty"`
}

type UserStore interface {
	Get(id uuid.UUID) (User, bool, error)
}
