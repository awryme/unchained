package trojanproto

type User struct {
	Desc string
	Key  Key
}

type UserStore interface {
	GetTrojan(key Key) (User, bool, error)
}
