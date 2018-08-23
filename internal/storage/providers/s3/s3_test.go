package s3

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
	"github.com/tuhlz/errors"
	"github.com/tuhlz/growtv/src/internal/storage"
)

var Session *session.Session
var Bucket *string

func init() {
	var err error

	err = gotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	Session, err = session.NewSession()
	if err != nil {
		panic(err.Error())
	}

	Bucket = aws.String(os.Getenv("GROWTV_BUCKET_TESTING"))
	if Bucket == nil {
		panic("Bucket was not set")
	}
}

func mockStore() storage.Storage {
	return New(&Config{
		Bucket:  Bucket,
		Session: Session,
	})
}

func TestNew(t *testing.T) {
	store := New(&Config{
		Bucket:  Bucket,
		Session: Session,
		ContentTypeDecoder: func(_ string) (contentType string) {
			return "text/plain"
		},
	})

	assert.NotNil(t, store)
}

type s3Tester struct {
	key  string
	typ  string
	data []byte
}

func (t *s3Tester) SetData(data []byte) {
	t.data = data
}

func (t *s3Tester) Key() string {
	return t.key
}

func (t *s3Tester) Type() string {
	if t.typ != "" {
		return t.typ
	}

	return "s3.s3Tester"
}

func (t *s3Tester) Data() []byte {
	return t.data
}

func TestDb_Write(t *testing.T) {
	store := mockStore()

	id := xid.New()

	a := s3Tester{key: id.String(), data: id.Bytes()}
	err := store.Write(&a)
	assert.Nil(t, err)
}

func TestDb_Read(t *testing.T) {
	store := mockStore()

	id := xid.New()

	a := s3Tester{key: id.String(), data: id.Bytes()}
	err := store.Write(&a)
	assert.Nil(t, err)

	b := s3Tester{key: a.key}
	err = store.Read(&b)
	assert.Nil(t, err)

	assert.Equal(t, a, b)
}

func TestDb_Delete(t *testing.T) {
	store := mockStore()

	id := xid.New()

	a := s3Tester{key: id.String(), data: id.Bytes()}
	err := store.Write(&a)
	assert.Nil(t, err)

	b := s3Tester{key: a.key}
	err = store.Read(&b)
	assert.Nil(t, err)

	assert.Equal(t, a, b)

	err = store.Delete(&a)
	assert.Nil(t, err)

	err = store.Read(&b)
	assert.NotNil(t, err)
	assert.Equal(t, 404, errors.Code(err), errors.Message(err))
}

type s3TestLister struct {
	typ   string
	count int64
	rows  []s3Tester
}

func (t *s3TestLister) Type() string {
	return t.typ
}

func (t *s3TestLister) SetCount(count int64) {
	t.count = count
}

func (t *s3TestLister) Append(id, typ string, data []byte) {
	t.rows = append(t.rows, s3Tester{key: id, typ: typ, data: data})
}

func TestDb_List(t *testing.T) {
	store := mockStore()

	groupA := xid.New().String()
	for range []int{1, 2, 3, 4} {
		id := xid.New()

		store.Write(&s3Tester{key: id.String(), typ: groupA, data: id.Bytes()})
	}

	groupB := xid.New().String()
	for range []int{1, 2, 3, 4} {
		id := xid.New()

		store.Write(&s3Tester{key: id.String(), typ: groupB, data: id.Bytes()})
	}

	a := s3TestLister{typ: groupB}

	err := store.List(&a)
	assert.Nil(t, err)
	assert.Equal(t, int(a.count), len(a.rows))
	assert.Equal(t, 4, int(a.count))
}

func TestDb_DecodeKey(t *testing.T) {
	d := &db{
		config: &Config{
			KeyEncoder: DefaultKeyEncoder,
			KeyDecoder: DefaultKeyDecoder,
		},
	}

	a, b := xid.New().String(), xid.New().String()

	key := d.config.KeyEncoder(a, b)
	id, typ := d.config.KeyDecoder(key)

	assert.Equal(t, a, id)
	assert.Equal(t, b, typ)
}
