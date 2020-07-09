package http_handler

import (
	"errors"
)

var (
	ErrEmptyBody       = errors.New("empty body")
	ErrBodyParseFailed = errors.New("body parse failed")
)
