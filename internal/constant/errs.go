package constant

import "errors"

const (
	CodeErrInternalServer = 500
	CodeSuccess           = 200
	CodeErrBadRequest     = 400
)

// internal
var (
	ErrRecordExist     = errors.New("ErrRecordExist")
	ErrRecordNotFound  = errors.New("ErrRecordNotFound")
	ErrInvalidParams   = errors.New("ErrInvalidParams")
	ErrStructTransform = errors.New("ErrStructTransform")
)

// api
var (
	ErrTypeInternalServer = "ErrInternalServer"
	ErrTypeInvalidParams  = "ErrInvalidParams"
)

// app
var (
	ErrCmdTimeout = "ErrCmdTimeout"
)
