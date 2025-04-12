package singleuserstore

import (
	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
)

type Trojan struct {
	user trojanproto.User
}

func NewTrojan(name string, pass string) *Trojan {
	return &Trojan{
		user: trojanproto.User{
			Name: name,
			Key:  trojanproto.NewKey(pass),
		},
	}
}

func (s *Trojan) Get(key trojanproto.Key) (trojanproto.User, bool, error) {
	if s.user.Key != key {
		return s.user, false, nil
	}
	return s.user, true, nil
}
