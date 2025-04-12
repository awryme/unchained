package users

import (
	"github.com/gofrs/uuid/v5"
)

type Vless struct {
	Name string    `json:"name"`
	UUID uuid.UUID `json:"uuid"`
	Flow string    `json:"flow,omitempty"`
}

type Trojan struct {
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}
