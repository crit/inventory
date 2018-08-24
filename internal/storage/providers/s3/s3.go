package s3

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/crit/inventory/internal/storage"
)

const Separator = "-|-"

type Config struct {
	Domain             string
	Bucket             *string
	Session            *session.Session
	Expires            *time.Time
	KeyDecoder         func(value string) (key, typ string)
	KeyEncoder         func(key, typ string) string
	ContentTypeDecoder func(key string) (contentType string)
}

func DefaultContentTypeDecoder(_ string) (contentType string) {
	return "application/octet-stream"
}

func DefaultKeyDecoder(value string) (key, typ string) {
	key = strings.TrimSuffix(value, filepath.Ext(value))
	parts := strings.Split(key, Separator)

	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return key, ""
}

func DefaultKeyEncoder(key, typ string) string {
	ext := filepath.Ext(key)
	k := strings.TrimSuffix(key, ext)

	if ext == "" {
		ext = ".bin"
	}

	return fmt.Sprintf("%s-|-%s%s", k, typ, ext)
}

func New(config *Config) storage.Storage {
	if config.KeyDecoder == nil {
		config.KeyDecoder = DefaultKeyDecoder
	}

	if config.KeyEncoder == nil {
		config.KeyEncoder = DefaultKeyEncoder
	}

	if config.ContentTypeDecoder == nil {
		config.ContentTypeDecoder = DefaultContentTypeDecoder
	}

	return &db{
		uploader:   s3manager.NewUploader(config.Session),
		downloader: s3manager.NewDownloader(config.Session),
		deleter:    s3manager.NewBatchDelete(config.Session),
		svc:        s3.New(config.Session),
		config:     config,
	}
}

type db struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	deleter    *s3manager.BatchDelete
	svc        *s3.S3
	config     *Config
}

func (d *db) Write(model storage.Writer) error {
	key := d.config.KeyEncoder(model.Key(), model.Type())
	typ := d.config.ContentTypeDecoder(model.Key())

	input := &s3manager.UploadInput{
		Bucket:      d.config.Bucket,
		ContentType: aws.String(typ),
		Expires:     d.config.Expires,
		Key:         aws.String(key),
		Body:        bytes.NewReader(model.Data()),
	}

	_, err := d.uploader.Upload(input)

	return toError(err)
}

func (d *db) Read(model storage.Reader) error {
	var buff aws.WriteAtBuffer

	key := d.config.KeyEncoder(model.Key(), model.Type())

	input := &s3.GetObjectInput{
		Bucket: d.config.Bucket,
		Key:    aws.String(key),
	}

	_, err := d.downloader.Download(&buff, input)

	if err != nil {
		return toError(err)
	}

	model.SetData(buff.Bytes())

	return nil
}

func (d *db) Delete(model storage.Meta) error {
	key := d.config.KeyEncoder(model.Key(), model.Type())

	input := &s3.DeleteObjectInput{
		Bucket: d.config.Bucket,
		Key:    aws.String(key),
	}

	_, err := d.svc.DeleteObject(input)

	if err != nil {
		return toError(err)
	}

	return nil
}

func (d *db) List(model storage.Mapper) error {
	var count int64

	input := &s3.ListObjectsInput{
		Bucket: d.config.Bucket,
	}

	err := d.svc.ListObjectsPages(input, func(out *s3.ListObjectsOutput, lastPage bool) bool {
		for _, obj := range out.Contents {
			id, typ := d.config.KeyDecoder(aws.StringValue(obj.Key))

			if model.Type() == typ {
				model.Append(id, typ, nil)
				count++
			}
		}

		return true
	})

	if err != nil {
		return toError(err)
	}

	model.SetCount(count)

	return nil
}
