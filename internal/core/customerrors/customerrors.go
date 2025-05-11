package customerrors

import "errors"

var (
	ErrIconNotFound    = errors.New("icon not found")
	ErrInvalidIconName = errors.New("invalid icon name")
)
