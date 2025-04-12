package trojanproto

type User struct {
	Name string
	Key  Key
}

type UserStore interface {
	Get(key Key) (User, bool, error)
}
