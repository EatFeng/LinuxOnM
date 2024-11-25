package copier

import (
	"encoding/json"
	"github.com/pkg/errors"
)

func Copy(dst, src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return errors.Wrap(err, "failed to marshal src")
	}
	if err = json.Unmarshal(b, dst); err != nil {
		return errors.Wrap(err, "failed to unmarshal src")
	}
	return nil
}
