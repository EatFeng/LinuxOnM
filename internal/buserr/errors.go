package buserr

import (
	"fmt"
	"github.com/pkg/errors"
)

type BusinessError struct {
	Msg    string
	Detail interface{}
	Map    map[string]interface{}
	Err    error
}

func (e BusinessError) Error() string {
	content := e.Msg
	if e.Detail != nil {
		content = fmt.Sprintf("%s: %v", e.Msg, e.Detail)
	} else if e.Map != nil {
		content = fmt.Sprintf("%s: %v", e.Msg, e.Map)
	}
	if content == "" {
		if e.Err != nil {
			return e.Err.Error()
		}
		return errors.New(e.Msg).Error()
	}
	return content
}

func New(Key string) BusinessError {
	return BusinessError{
		Msg:    Key,
		Detail: nil,
		Err:    nil,
	}
}

func WithDetail(Key string, detail interface{}, err error) BusinessError {
	return BusinessError{
		Msg:    Key,
		Detail: detail,
		Err:    err,
	}
}

func WithName(Key string, name string) BusinessError {
	paramMap := map[string]interface{}{}
	if name != "" {
		paramMap["name"] = name
	}
	return BusinessError{
		Msg: Key,
		Map: paramMap,
	}
}
