package constant

import "errors"

const (
	CodeErrInternalServer = 500
	CodeAuth              = 406
	CodeErrNotFound       = 404
	CodeErrUnauthorized   = 401
	CodeErrBadRequest     = 400
	CodePasswordExpired   = 313
	CodeErrEntrance       = 312
	CodeErrIP             = 310
	CodeSuccess           = 200
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
	ErrTokenParse      = errors.New("ErrTokenParse")
)

// api
var (
	ErrTypeInternalServer  = "ErrInternalServer"
	ErrTypeInvalidParams   = "ErrInvalidParams"
	ErrTypePasswordExpired = "ErrPasswordExpired"
	ErrCmdIllegal          = "ErrCmdIllegal"
	ErrNameIsExist         = "ErrNameIsExist"
	ErrTypeNotLogin        = "ErrNotLogin"
)

// app
var (
	ErrCmdTimeout     = "ErrCmdTimeout"
	ErrFileCanNotRead = "ErrFileCanNotRead"
	ErrPortInUsed     = "ErrPortInUsed"
	ErrFileNotFound   = "ErrFileNotFound"
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
	ErrMovePathFailed   = "ErrMovePathFailed"
)

// firewall
var (
	ErrFirewall = "ErrFirewall"
)
