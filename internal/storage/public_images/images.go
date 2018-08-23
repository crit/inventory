package public_images

import (
	"fmt"
	"sync"

	"github.com/tuhlz/growtv/src/internal/storage"
	"github.com/tuhlz/growtv/src/internal/storage/providers"
	"github.com/tuhlz/growtv/src/internal/storage/providers/s3"
)

var (
	db     = providers.Mock(nil)
	config *s3.Config

	lock sync.RWMutex
)

func Register(store storage.Storage, cfg *s3.Config) {
	lock.Lock()
	defer lock.Unlock()

	db = store
	config = cfg
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

func EncodeURL(model storage.Meta) string {
	lock.RLock()
	defer lock.RUnlock()

	key := config.KeyEncoder(model.Key(), model.Type())

	return fmt.Sprintf("%s/%s", config.Domain, key)
}

func DecodeURL(value string) (key string) {
	lock.RLock()
	defer lock.RUnlock()

	key, _ = config.KeyDecoder(value)
	return key
}
