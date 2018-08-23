package dynamodb

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/tuhlz/errors"
)

func toError(err error) error {
	e, ok := err.(awserr.Error)

	if !ok {
		return errors.New(http.StatusInternalServerError, err)
	}

	if strings.Contains(e.Code(), "not found") {
		return errors.String(http.StatusNotFound, e.Message())
	}

	return errors.String(http.StatusInternalServerError, e.Message())
}
