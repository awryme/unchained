package userstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"os"

	"github.com/awryme/subdir"
	"github.com/awryme/unchained/pkg/protocols/trojan/trojanproto"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessproto"
	"github.com/awryme/unchained/pkg/protocols/vless/vlessvision"
	"github.com/gofrs/uuid/v5"
	"github.com/sethvargo/go-password/password"
)

var ErrUserNotFound = fmt.Errorf("user not found")

const fileProxyUsers = "proxy_users.json"

type UserInfo struct {
	ID             uuid.UUID
	Desc           string
	VlessID        uuid.UUID
	TrojanPassword string
}

type Users struct {
	Users map[uuid.UUID]UserInfo
}

type FileStore struct {
	dir subdir.Dir

	currentUsers     Users
	usersByVlessID   map[uuid.UUID]UserInfo
	usersByTrojanKey map[trojanproto.Key]UserInfo
}

func NewFileStore(dir string) (*FileStore, error) {
	subDir, err := subdir.New(dir, os.ModeDir)
	if err != nil {
		return nil, fmt.Errorf("create new subdir: %w", err)
	}
	store := &FileStore{
		dir: subDir,
	}
	err = store.readUsersFile()
	if err != nil {
		return nil, fmt.Errorf("read users file: %w", err)
	}
	return store, nil
}

func (store *FileStore) GetUser(id uuid.UUID) (UserInfo, error) {
	user, ok := store.currentUsers.Users[id]
	if !ok {
		return user, fmt.Errorf("get user %s: %w", id.String(), ErrUserNotFound)
	}
	return user, nil
}

func (store *FileStore) Add(id uuid.UUID, desc string) error {
	user := UserInfo{
		ID:   id,
		Desc: desc,
	}
	err := genUser(&user)
	if err != nil {
		return fmt.Errorf("generate new user: %w", err)
	}

	store.currentUsers.Users[user.ID] = user
	err = store.writeUsersFile()
	if err != nil {
		return fmt.Errorf("write users file: %w", err)
	}

	store.usersByVlessID[user.VlessID] = user

	key := trojanproto.NewKey(user.TrojanPassword)
	store.usersByTrojanKey[key] = user

	return nil
}

func (store *FileStore) RemoveUser(id uuid.UUID) error {
	user, ok := store.currentUsers.Users[id]
	if !ok {
		return fmt.Errorf("delete user %s: %w", id.String(), ErrUserNotFound)
	}
	delete(store.currentUsers.Users, user.ID)
	delete(store.usersByVlessID, user.VlessID)
	delete(store.usersByTrojanKey, trojanproto.NewKey(user.TrojanPassword))

	return store.writeUsersFile()
}

func (store *FileStore) ListIDs() iter.Seq[uuid.UUID] {
	return func(yield func(uuid.UUID) bool) {
		for id := range store.currentUsers.Users {
			if !yield(id) {
				return
			}
		}
	}
}

// accepts vless user id
func (store *FileStore) GetVless(vlessID uuid.UUID) (vlessproto.User, bool, error) {
	var vlessUser vlessproto.User
	user, ok := store.usersByVlessID[vlessID]
	if !ok {
		return vlessUser, false, nil
	}

	vlessUser = vlessproto.User{
		UUID: user.VlessID,
		Desc: user.Desc,
		Flow: vlessvision.Flow,
	}

	return vlessUser, true, nil
}

func (store *FileStore) GetTrojan(key trojanproto.Key) (trojanproto.User, bool, error) {
	var trojanUser trojanproto.User
	user, ok := store.usersByTrojanKey[key]
	if !ok {
		return trojanUser, false, nil
	}

	trojanUser = trojanproto.User{
		Desc: user.Desc,
		Key:  key,
	}

	return trojanUser, true, nil
}

func (store *FileStore) readUsersFile() error {
	store.usersByVlessID = make(map[uuid.UUID]UserInfo)
	store.usersByTrojanKey = make(map[trojanproto.Key]UserInfo)
	store.currentUsers = Users{
		Users: make(map[uuid.UUID]UserInfo),
	}

	file, err := store.dir.Open(fileProxyUsers)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read users file: %w", err)
	}
	defer file.Close()

	var users Users
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return fmt.Errorf("decode users file: %w", err)
	}
	store.currentUsers = users
	for _, user := range users.Users {
		store.usersByVlessID[user.VlessID] = user

		key := trojanproto.NewKey(user.TrojanPassword)
		store.usersByTrojanKey[key] = user
	}
	return nil
}

func (store *FileStore) writeUsersFile() error {
	file, err := store.dir.Create(fileProxyUsers)
	if err != nil {
		return fmt.Errorf("create users file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(store.currentUsers)
	if err != nil {
		return fmt.Errorf("encode users file: %w", err)
	}
	return nil
}

func genUser(user *UserInfo) error {
	const length = 16

	pwd, err := password.Generate(length, length/3, 0, false, false)
	if err != nil {
		return fmt.Errorf("generate trojan password: %w", err)
	}
	user.TrojanPassword = pwd

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate vless user uuid: %w", err)
	}
	user.VlessID = id

	return nil
}
