package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error is the base error structure.
type Error struct {
	Code    int
	Message string
}

// Error returns error code and message in a formatted string. Satisfies
// the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("[%d] - %s", e.Code, e.Message)
}

// Code attempts to cast the error to Error and determines the code. Nil is
// 200; Non-Error is 500.
func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}

	e, ok := err.(Error)

	if !ok {
		return http.StatusInternalServerError
	}

	return e.Code
}

// Message attempts to cast the error to Error and determines the message. Nil is an empty
// string; Non-Error is the error interface output.
func Message(err error) string {
	if err == nil {
		return ""
	}

	e, ok := err.(Error)

	if !ok {
		return err.Error()
	}

	return e.Message
}

// Error is a convenience function that creates a new Error with code and
// the Message of the error.
func New(code int, err error) error {
	return Error{code, Message(err)}
}

// ErrorString is a convenience function that creates a new Error with code and a message
func String(code int, msg string) error {
	return Error{code, msg}
}

// ErrorPackage is a convenience function to output a standard response for
// an HTTP handler of our specific API Service.
func Package(err error) map[string]interface{} {
	if err == nil {
		return nil
	}

	return map[string]interface{}{
		"code":    Code(err),
		"message": Message(err),
	}
}

func UnPack(raw []byte) Error {
	var obj map[string]interface{}
	var out Error
	var success bool

	err := json.Unmarshal(raw, &obj)

	if err != nil {
		out.Code = 500
		out.Message = string(raw)
		return out
	}

	if t, ok := obj["code"]; ok {
		out.Code, success = t.(int)
		if !success {
			out.Code = 500
		}
	}

	if t, ok := obj["message"]; ok {
		out.Message, success = t.(string)
		if !success {
			out.Message = string(raw)
		}
	}

	return out
}

// Affected is a convenience function for handling whether a db update touched the correct
// number of records. Usually causes the update to be rolled back if error is not nil.
func Affected(expected, actual int64) error {
	if expected == actual {
		return nil
	}

	return New(
		http.StatusInternalServerError,
		fmt.Errorf("expected %d affected rows; got %d rows; no update made", expected, actual),
	)
}
