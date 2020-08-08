package upngo

import (
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

type ErrorObject struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type SourceObject struct {
	Parameter string `json:"parameter"`
	Pointer   string `json:"pointer"`
}

type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
	Source SourceObject  `json:"source"`
}
