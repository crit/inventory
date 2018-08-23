package s3

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tuhlz/errors"
)

func toError(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(awserr.Error)

	if !ok {
		return errors.New(http.StatusInternalServerError, err)
	}

	switch e.Code() {
	case s3.ErrCodeNoSuchKey, s3.ErrCodeNoSuchUpload, s3.ErrCodeNoSuchBucket:
		return errors.String(http.StatusNotFound, e.Message())
	case s3.ErrCodeBucketAlreadyExists:
		return nil
	}

	return errors.String(http.StatusInternalServerError, e.Message())
}
