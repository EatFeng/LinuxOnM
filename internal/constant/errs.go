package constant

import "errors"

const (
	CodeErrInternalServer = 500
	CodeSuccess           = 200
	CodeErrBadRequest     = 400
	CodeAuth              = 406
	CodePasswordExpired   = 313
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
	ErrNotSupportType  = errors.New("ErrNotSupportType")
)

// api
var (
	ErrTypeInternalServer  = "ErrInternalServer"
	ErrTypeInvalidParams   = "ErrInvalidParams"
	ErrTypePasswordExpired = "ErrPasswordExpired"
)

// app
var (
	ErrCmdTimeout     = "ErrCmdTimeout"
	ErrFileCanNotRead = "ErrFileCanNotRead"
)

var (
	ErrEntrance    = "ErrEntrance"
	ErrGroupIsUsed = "ErrGroupIsUsed"
)

// file
var (
	ErrLinkPathNotFound = "ErrLinkPathNotFound"
	ErrFileIsExist      = "ErrFileIsExist"
)
