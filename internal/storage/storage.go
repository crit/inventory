package storage

import "github.com/rs/xid"

type Item interface {
	Key() string
}

type Group interface {
	Type() string
}

type Meta interface {
	Item
	Group
}

type Writer interface {
	Item
	Group
	Data() []byte
}

type Reader interface {
	Item
	Group
	SetData(data []byte)
}

type Mapper interface {
	Group
	SetCount(count int64)
	Append(id, typ string, data []byte)
}

type Storage interface {
	Write(model Writer) error
	Read(model Reader) error
	Delete(model Meta) error
	List(model Mapper) error
}

type TypeDecoder func(key string) (ext string, contentType string)
type KeyDecoder func(value string) (key, typ string)
type KeyEncoder func(key, typ string) string

func NextID() string {
	return xid.New().String()
}
