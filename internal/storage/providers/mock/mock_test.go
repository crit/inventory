package mock

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	store := New(func(db *DB) {
		assert.NotNil(t, db)
	})

	assert.NotNil(t, store)
}

type mockTester struct {
	key  string
	typ  string
	data []byte
}

func (m *mockTester) SetData(data []byte) {
	m.data = data
}

func (m *mockTester) Key() string {
	return m.key
}

func (m *mockTester) Type() string {
	return m.typ
}

func (m *mockTester) Data() []byte {
	return m.data
}

func TestMockDB_Write(t *testing.T) {
	var mock *DB
	store := New(func(instance *DB) {
		mock = instance
	})

	a := mockTester{key: "test-key", typ: "storage.mockTester", data: []byte("test-value")}
	err := store.Write(&a)
	assert.Nil(t, err)
	assert.Len(t, mock.Store, 1)

	v := mock.Store[mock.Key(&a)]
	assert.Equal(t, a.data, v)

	for range []int{1, 2, 3, 4} {
		err := store.Write(&mockTester{
			key:  xid.New().String(),
			typ:  "storage.mockTester",
			data: xid.New().Bytes(),
		})
		assert.Nil(t, err)
	}
}

func TestMockDB_Read(t *testing.T) {
	store := New(nil)

	a := mockTester{key: "test-key", typ: "storage.mockTester", data: []byte("test-value")}
	err := store.Write(&a)
	assert.Nil(t, err)

	b := mockTester{key: a.key, typ: a.typ}
	err = store.Read(&b)
	assert.Nil(t, err)
	assert.Equal(t, a, b)
}

func TestMockDB_Delete(t *testing.T) {
	store := New(nil)

	a := mockTester{key: "test-key", typ: "storage.mockTester", data: []byte("test-value")}
	err := store.Write(&a)
	assert.Nil(t, err)

	b := mockTester{key: a.key, typ: a.typ}
	err = store.Delete(&b)
	assert.Nil(t, err)

	c := mockTester{key: a.key, typ: a.typ}
	err = store.Read(&c)
	assert.NotNil(t, err)
	assert.Empty(t, c.data)
}

type mockTestLister struct {
	count int64
	items []mockTester
}

func (m *mockTestLister) Type() string {
	return "storage.mockTester"
}

func (m *mockTestLister) SetCount(count int64) {
	m.count = count
}

func (m *mockTestLister) Append(id, typ string, data []byte) {
	m.items = append(m.items, mockTester{
		key:  id,
		typ:  typ,
		data: data,
	})
}

func TestMockDB_List(t *testing.T) {
	var mock *DB
	store := New(func(db *DB) {
		mock = db
	})

	for range []int{1, 2, 3, 4} {
		id := xid.New()

		store.Write(&mockTester{
			key:  id.String(),
			typ:  "storage.mockTester",
			data: id.Bytes(),
		})
	}

	for range []int{1, 2, 3, 4} {
		id := xid.New()

		store.Write(&mockTester{
			key:  id.String(),
			typ:  "Unknown",
			data: id.Bytes(),
		})
	}

	assert.Len(t, mock.Store, 8)

	var b mockTestLister

	err := store.List(&b)
	assert.Nil(t, err)
	assert.Equal(t, b.count, int64(4))
}
