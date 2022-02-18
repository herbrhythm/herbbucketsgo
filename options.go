package herbbucketsgo

import (
	"net/http"
	"net/url"
)

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

func (o *CompleteOptions) Encode() *url.Values {
	params := &url.Values{}
	params.Add("id", o.ID)
	params.Add("bucket", o.Bucket)
	params.Add("objcet", o.Object)
	params.Add("ts", o.Timestamp)
	params.Add("sign", o.Sign)
	return params
}
func (o *CompleteOptions) Decode(r *http.Request) {
	q := r.URL.Query()
	o.ID = q.Get("id")
	o.Bucket = q.Get("bucket")
	o.Object = q.Get("object")
	o.Timestamp = q.Get("ts")
	o.Sign = q.Get("sign")
}

type CompleteOptions struct {
	ID        string
	Bucket    string
	Object    string
	Timestamp string
	Sign      string
}
