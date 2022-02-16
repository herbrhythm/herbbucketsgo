package herbbucketsgo

import (
	"errors"
	"net/http"

	"github.com/herb-go/fetcher"
)

func convertErr(err error) error {
	if fetcher.CompareResponseErrStatusCode(err, 404) {
		return ErrNotExist
	}
	if fetcher.CompareResponseErrStatusCode(err, 422) {
		return &UnprocessableEntityError{
			Response: err.(*fetcher.Response),
		}
	}
	return err
}

var ErrNotExist = errors.New("herbbucketsgo:bucket object not found")

type UnprocessableEntityError struct {
	*fetcher.Response
}

func (e *UnprocessableEntityError) MustDump(w http.ResponseWriter) {
	w.WriteHeader(422)
	data, err := e.Response.BodyContent()
	if err != nil {
		panic(err)
	}
	_, err = w.Write(data)
	panic(err)
}
