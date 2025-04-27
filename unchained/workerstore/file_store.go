package workerstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/awryme/subdir"
	"github.com/awryme/unchained-control/unchained-control/worker"
	"github.com/awryme/unchained/unchained/protocols"
	"github.com/awryme/unchained/unchained/urlmaker"
	"github.com/awryme/unchained/unchained/userstore"
	"github.com/gofrs/uuid/v5"
)

const fileWorkerData = "worker_data.json"

// FileStore implements worker.ApiStore
type FileStore struct {
	dir       subdir.Dir
	userstore *userstore.FileStore
	urlMaker  *urlmaker.UrlMaker
}

func NewFileStore(dir string, userstore *userstore.FileStore, urlMaker *urlmaker.UrlMaker) (*FileStore, error) {
	subDir, err := subdir.New(dir, os.ModeDir)
	if err != nil {
		return nil, fmt.Errorf("create new subdir: %w", err)
	}
	return &FileStore{subDir, userstore, urlMaker}, nil
}

// manage secrets and data
func (store *FileStore) GetData() (worker.WorkerData, bool, error) {
	var data worker.WorkerData

	file, err := store.dir.Open(fileWorkerData)
	if errors.Is(err, os.ErrNotExist) {
		return data, false, nil
	}
	if err != nil {
		return data, false, fmt.Errorf("create worker data file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return data, false, fmt.Errorf("decode worker data: %w", err)
	}
	return data, true, nil
}

func (store *FileStore) StoreData(secrets worker.WorkerData) error {
	file, err := store.dir.Create(fileWorkerData)
	if err != nil {
		return fmt.Errorf("create worker data file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(secrets)
	if err != nil {
		return fmt.Errorf("encode worker data file: %w", err)
	}
	return nil
}

// manage users
func (store *FileStore) SyncUsers(users []worker.UserInfo) error {
	currentIDs := make(map[uuid.UUID]bool)
	for id := range store.userstore.ListIDs() {
		currentIDs[id] = true
	}

	toAdd := make([]worker.UserInfo, 0)

	for _, user := range users {
		if currentIDs[user.ID] {
			delete(currentIDs, user.ID)
		} else {
			toAdd = append(toAdd, user)
		}
	}

	for _, user := range toAdd {
		err := store.userstore.Add(user.ID, user.Desc)
		if err != nil {
			return fmt.Errorf("add user to store: %w", err)
		}
	}

	// iterate by remaining to delete
	for id := range currentIDs {
		err := store.userstore.RemoveUser(id)
		if err != nil {
			return fmt.Errorf("remove user from store: %w", err)
		}
	}

	return nil
}

// list proxies for user
func (store *FileStore) GetUserProxyURLs(id uuid.UUID) ([]worker.ProxyURL, error) {
	user, err := store.userstore.GetUser(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	res := make([]worker.ProxyURL, 0, 2)

	res = append(res, worker.ProxyURL{
		Proto: protocols.Vless,
		URL:   store.urlMaker.MakeVless(user.VlessID),
	})

	res = append(res, worker.ProxyURL{
		Proto: protocols.Trojan,
		URL:   store.urlMaker.MakeTrojan(user.TrojanPassword),
	})

	return res, nil
}
