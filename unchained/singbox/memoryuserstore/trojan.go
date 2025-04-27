package memoryuserstore

import (
	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
)

type Trojan struct {
	users map[trojanproto.Key]trojanproto.User
}

func NewTrojan() *Trojan {
	return &Trojan{
		users: make(map[trojanproto.Key]trojanproto.User),
	}
}

func (s *Trojan) Add(desc, password string) error {
	key := trojanproto.NewKey(password)
	s.users[key] = trojanproto.User{
		Desc: desc,
		Key:  key,
	}
	return nil
}

func (s *Trojan) GetTrojan(key trojanproto.Key) (trojanproto.User, bool, error) {
	user, ok := s.users[key]
	return user, ok, nil
}
