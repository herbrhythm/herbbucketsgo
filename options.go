package herbbucketsgo

import "net/http"

type UploadOptions struct {
	TTL       string
	Filename  string
	SizeLimit string
	Size      string
}

func MustLoadUploadOptions(r *http.Request) *UploadOptions {
	q := r.URL.Query()
	opt := &UploadOptions{
		TTL:       q.Get(QueryFieldTTL),
		Filename:  q.Get(QueryFieldFilename),
		Size:      q.Get(QueryFieldSize),
		SizeLimit: q.Get(QueryFieldSizeLimit),
	}
	return opt
}
