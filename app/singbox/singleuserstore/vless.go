package singleuserstore

import (
	"github.com/awryme/unchained/pkg/protocols/vless/vlessproto"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessvision"
	"github.com/gofrs/uuid/v5"
)

type Vless struct {
	user vlessproto.User
}

func NewVless(name string, id uuid.UUID) *Vless {
	return &Vless{
		user: vlessproto.User{
			Name: name,
			UUID: id,
			Flow: vlessvision.Flow,
		},
	}
}

func (s *Vless) Get(id uuid.UUID) (vlessproto.User, bool, error) {
	if s.user.UUID != id {
		return s.user, false, nil
	}
	return s.user, true, nil
}
