package constant

import "errors"

const (
	CodeErrInternalServer = 500
	CodeSuccess           = 200
	CodeErrBadRequest     = 400
	CodeAuth              = 406
	CodePasswordExpired   = 313
	CodeErrNotFound       = 404
	CodeErrEntrance       = 312
	CodeErrIP             = 310
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
	ErrCmdIllegal          = "ErrCmdIllegal"
)

// app
var (
	ErrCmdTimeout     = "ErrCmdTimeout"
	ErrFileCanNotRead = "ErrFileCanNotRead"
	ErrPortInUsed     = "ErrPortInUsed"
)

var (
	ErrEntrance    = "ErrEntrance"
	ErrGroupIsUsed = "ErrGroupIsUsed"
)

// file
var (
	ErrLinkPathNotFound = "ErrLinkPathNotFound"
	ErrFileIsExist      = "ErrFileIsExist"
	ErrPathNotDelete    = "ErrPathNotDelete"
	ErrPathNotFound     = "ErrPathNotFound"
)

// firewall
var (
	ErrFirewall = "ErrFirewall"
)
