package bolt

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/tuhlz/growtv/src/internal/storage"
)

func boltDbBuilder() (storage.Storage, func()) {
	f, err := ioutil.TempFile("", "")

	if err != nil {
		panic(err.Error())
	}

	db, err := New(f.Name())

	if err != nil {
		panic(err.Error())
	}

	return db, func() {
		os.Remove(f.Name())
	}
}

type boltTester struct {
	id   string
	typ  string
	data []byte
}

func (b *boltTester) SetData(data []byte) {
	b.data = data
}

func (b *boltTester) Key() string {
	return b.id
}

func (b *boltTester) Type() string {
	if b.typ == "" {
		return "storage.boltTester"
	}

	return b.typ
}

func (b *boltTester) Data() []byte {
	return b.data
}

func TestBoltDB_Write(t *testing.T) {
	store, fn := boltDbBuilder()
	defer fn()

	id := xid.New()
	err := store.Write(&boltTester{id: id.String(), data: id.Bytes()})
	assert.Nil(t, err)
}

func TestBoltDB_Read(t *testing.T) {
	store, fn := boltDbBuilder()
	defer fn()

	id := xid.New()
	a := boltTester{id: id.String(), data: id.Bytes()}
	err := store.Write(&a)
	assert.Nil(t, err)

	b := boltTester{id: a.id}
	err = store.Read(&b)
	assert.Nil(t, err)
	assert.Equal(t, a, b)
}

func TestBoltDB_Delete(t *testing.T) {
	store, fn := boltDbBuilder()
	defer fn()

	id := xid.New()
	a := boltTester{id: id.String(), data: id.Bytes()}
	err := store.Write(&a)
	assert.Nil(t, err)

	b := boltTester{id: a.id}
	err = store.Delete(&b)
	assert.Nil(t, err)

	c := boltTester{id: a.id}
	err = store.Read(&c)
	assert.NotNil(t, err)
}

type boltTestLister struct {
	count int64
	items []boltTester
}

func (b *boltTestLister) Type() string {
	return "storage.boltTester"
}

func (b *boltTestLister) SetCount(count int64) {
	b.count = count
}

func (b *boltTestLister) Append(id, typ string, data []byte) {
	b.items = append(b.items, boltTester{
		id:   id,
		typ:  typ,
		data: data,
	})
}

func TestBoltDB_List(t *testing.T) {
	store, fn := boltDbBuilder()
	defer fn()

	for range []int{1, 2, 3, 4} {
		id := xid.New()

		err := store.Write(&boltTester{
			id:   id.String(),
			typ:  "storage.boltTester",
			data: id.Bytes(),
		})

		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	for range []int{1, 2, 3, 4} {
		id := xid.New()

		err := store.Write(&boltTester{
			id:   id.String(),
			typ:  "unknown",
			data: id.Bytes(),
		})

		if err != nil {
			assert.Fail(t, err.Error())
		}
	}

	a := boltTestLister{}

	err := store.List(&a)
	assert.Nil(t, err)
	assert.Equal(t, len(a.items), int(a.count))
	assert.Equal(t, 4, int(a.count))
}
