package mock

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/tuhlz/errors"
	"github.com/tuhlz/growtv/src/internal/storage"
)

func New(reporter func(db *DB)) storage.Storage {
	store := &DB{
		Store: map[string][]byte{},
	}

	if reporter != nil {
		reporter(store)
	}

	return store
}

type typeRegistry struct {
	Type string
	Key  string
}

type DB struct {
	Store map[string][]byte

	types    []typeRegistry
	reporter func(mock *DB)
	lock     sync.RWMutex
}

func (DB) Key(model storage.Meta) string {
	return fmt.Sprintf("%s-%s", model.Key(), model.Type())
}

func (db *DB) Write(model storage.Writer) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	key := db.Key(model)

	db.Store[key] = model.Data()

	for _, registry := range db.types {
		if registry.Key == key {
			return nil // already indexed
		}
	}

	db.types = append(db.types, typeRegistry{
		Type: model.Type(),
		Key:  key,
	})

	return nil
}

func (db *DB) Read(model storage.Reader) error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	data, ok := db.Store[db.Key(model)]

	if !ok {
		return errors.String(http.StatusNotFound, "not found")
	}

	model.SetData(data)

	return nil
}

func (db *DB) Delete(model storage.Meta) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	delete(db.Store, db.Key(model))

	return nil
}

func (db *DB) List(mapper storage.Mapper) error {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var count int64

	for _, registry := range db.types {
		if registry.Type == mapper.Type() {
			mapper.Append(registry.Key, registry.Type, db.Store[registry.Key])
			count++
		}
	}

	mapper.SetCount(count)

	return nil
}
