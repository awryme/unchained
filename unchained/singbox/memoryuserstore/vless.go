package memoryuserstore

import (
	"github.com/awryme/unchained/pkg/protocols/vless/vlessproto"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessvision"
	"github.com/gofrs/uuid/v5"
)

type Vless struct {
	users map[uuid.UUID]vlessproto.User
}

func NewVless() *Vless {
	return &Vless{
		users: make(map[uuid.UUID]vlessproto.User),
	}
}

func (s *Vless) Add(desc string, id uuid.UUID) error {
	s.users[id] = vlessproto.User{
		Desc: desc,
		UUID: id,
		Flow: vlessvision.Flow,
	}
	return nil
}

func (s *Vless) GetVless(id uuid.UUID) (vlessproto.User, bool, error) {
	u, ok := s.users[id]
	return u, ok, nil
}
