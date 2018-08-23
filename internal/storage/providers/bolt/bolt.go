package bolt

import (
	"fmt"
	"strings"

	"github.com/coreos/bbolt"
	"github.com/tuhlz/errors"
	"github.com/tuhlz/growtv/src/internal/storage"
)

func New(path string) (storage.Storage, error) {
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		return nil, errors.New(500, err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(path))
		return err
	})

	if err != nil {
		return nil, errors.New(500, err)
	}

	return &boltDB{
		path: []byte(path),
		svc:  db,
	}, nil
}

type boltDB struct {
	svc  *bolt.DB
	path []byte
}

func (boltDB) Key(model storage.Meta) []byte {
	return []byte(fmt.Sprintf("%s-|-%s", model.Key(), model.Type()))
}

func (boltDB) DecodeKey(key []byte) (id, typ string) {
	parts := strings.Split(string(key), "-|-")

	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return "", ""
}

func (db *boltDB) Write(model storage.Writer) error {
	return db.svc.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(db.path).Put(db.Key(model), model.Data())
	})
}

func (db *boltDB) Read(model storage.Reader) error {
	var value []byte

	db.svc.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(db.path).Get(db.Key(model))
		return nil
	})

	if value == nil {
		return errors.String(404, "not found")
	}

	model.SetData(value)

	return nil
}

func (db *boltDB) Delete(model storage.Meta) error {
	db.svc.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(db.path).Delete(db.Key(model))
	})

	return nil
}

func (db *boltDB) List(model storage.Mapper) error {
	var count int64

	db.svc.View(func(tx *bolt.Tx) error {
		tx.Bucket(db.path).ForEach(func(k, data []byte) error {
			id, typ := db.DecodeKey(k)

			if typ == model.Type() {
				model.Append(id, typ, data)
				count++
			}

			return nil
		})

		return nil
	})

	model.SetCount(count)

	return nil
}
