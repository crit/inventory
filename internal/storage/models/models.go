package models

import (
	"sync"
	"time"

	"encoding/json"

	"github.com/crit/inventory/internal/storage"
	"github.com/crit/inventory/internal/storage/providers"
	"github.com/rs/xid"
)

var (
	db = providers.Mock(nil)

	lock sync.RWMutex
)

type Model interface {
	Read() error
	Write() error
	Delete() error
}

type ModelList interface {
	Read() error
}

func Register(store storage.Storage) {
	lock.Lock()
	defer lock.Unlock()
	db = store
}

func Write(model storage.Writer) error {
	lock.RLock()
	defer lock.RUnlock()

	return db.Write(model)
}

func Read(model storage.Reader) error {
	lock.RLock()
	defer lock.RUnlock()

	return db.Read(model)
}

func Delete(model storage.Meta) error {
	lock.RLock()
	defer lock.RUnlock()

	return db.Delete(model)
}

func List(model storage.Mapper) error {
	lock.RLock()
	defer lock.RUnlock()

	return db.List(model)
}

type ListMeta struct {
	Count int64 `json:"count"`
}

type Tracking struct {
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToBytes(obj interface{}) []byte {
	out, _ := json.Marshal(obj)
	return out
}

func FromBytes(data []byte, obj interface{}) {
	json.Unmarshal(data, obj)
}

func NextID() string {
	return xid.New().String()
}
