package storage

import "errors"

var (
	ErrorModelNotFound             = errors.New("model not found")
	ErrorUniqueConstraintViolation = errors.New("model should be unique")
)
