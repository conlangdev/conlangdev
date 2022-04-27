package conlangdev

import (
	"fmt"
)

type Error struct {
	Code       string
	Message    string
	StatusCode int
}

type FieldsError struct {
	Code       string
	Message    string
	StatusCode int
	Fields     []string
}

const (
	ECONFLICT     = "conflict"
	ESERVER       = "server_error"
	EBADREQUEST   = "bad_request"
	EVALIDFAIL    = "validation_failed"
	ENOTFOUND     = "not_found"
	EUNAUTHORIZED = "unauthorized"
)

func (e *Error) Error() string {
	return fmt.Sprintf("conlangdev error (%s): %s", e.Code, e.Message)
}

func (e *FieldsError) Error() string {
	return fmt.Sprintf("conlangdev fields error (%s): %s %v", e.Code, e.Message, e.Fields)
}
