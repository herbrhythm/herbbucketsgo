package herbbucketsgo

import (
	"errors"

	"github.com/herb-go/fetcher"
)

func convertErr(err error) error {
	if fetcher.IsResponseErr(err) {
		if err.(*fetcher.Response).StatusCode == 404 {
			return ErrNotExist
		}
	}
	return err
}

var ErrNotExist = errors.New("herbbucketsgo:bucket object not found")
