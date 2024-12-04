package constant

import "errors"

const (
	CodeErrInternalServer = 500
	CodeSuccess           = 200
	CodeErrBadRequest     = 400
	CodeAuth              = 406
)

// internal
var (
	ErrRecordExist     = errors.New("ErrRecordExist")
	ErrRecordNotFound  = errors.New("ErrRecordNotFound")
	ErrInvalidParams   = errors.New("ErrInvalidParams")
	ErrStructTransform = errors.New("ErrStructTransform")
	ErrCaptchaCode     = errors.New("ErrCaptchaCode")
	ErrAuth            = errors.New("ErrAuth")
	ErrInitialPassword = errors.New("ErrInitialPassword")
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

var (
	ErrEntrance    = "ErrEntrance"
	ErrGroupIsUsed = "ErrGroupIsUsed"
)
