package herbbucketsgo

import (
	"io"
	"net/url"
	"strconv"

	"github.com/herb-go/fetcher"
)

const UploadTypePut = "put"
const UploadTypePost = "post"

type SaveResult struct {
	ID     string
	Bucket string
	Object string
}

type WebuploadInfo struct {
	ID             string
	Bucket         string
	Object         string
	UploadURL      string
	Complete       string
	UploadType     string
	SizeLimit      int64
	ExpiredAt      int64
	PostBody       map[string]string
	FileField      string
	SuccessCodeMin int
	SuccessCodeMax int
}

type DownloadInfo struct {
	URL       string
	ExpiredAt int64
	Permanent bool
}
type Fileinfo struct {
	Size    int64
	Modtime int64
}

func New() *Buckets {
	return &Buckets{}
}

type CompleteInfo struct {
	ID      string
	Bucket  string
	Object  string
	Size    int64
	Preview *DownloadInfo
}

type Buckets struct {
	Bucket                 string
	Passthrough            bool
	PresetGrantUploadInfo  *fetcher.Preset
	PresetGrantDownladInfo *fetcher.Preset
	PresetComplete         *fetcher.Preset
	PresetContent          *fetcher.Preset
	PresetSave             *fetcher.Preset
	PresetRemove           *fetcher.Preset
	PresetInfo             *fetcher.Preset
}

func (s *Buckets) GrantUploadInfo(bucket string, opt *UploadOptions) (*WebuploadInfo, error) {
	result := &WebuploadInfo{}
	params := url.Values{}
	params.Add(QueryFieldFilename, opt.Filename)
	params.Add(QueryFieldTTL, opt.TTL)
	params.Add(QueryFieldSizeLimit, opt.SizeLimit)
	params.Add(QueryFieldSize, opt.Size)
	_, err := s.PresetGrantUploadInfo.Concat(fetcher.PathJoin(bucket), fetcher.Params(params)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, convertErr(err)
	}
	return result, nil
}

func (s *Buckets) GrantDownloadInfo(bucket string, object string, ttl int64) (*DownloadInfo, error) {
	var result = &DownloadInfo{}
	params := url.Values{}
	params.Add(QueryFieldTTL, strconv.FormatInt(ttl, 10))
	_, err := s.PresetGrantDownladInfo.Concat(fetcher.PathJoin(bucket), fetcher.PathJoin(object), fetcher.Params(params)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, convertErr(err)
	}
	return result, nil
}

func (s *Buckets) Complete(opt *CompleteOptions) (*CompleteInfo, error) {
	var result = &CompleteInfo{}
	_, err := s.PresetComplete.Concat(fetcher.Params(*opt.Encode())).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, convertErr(err)
	}
	return result, nil
}
func (s *Buckets) contentPassthrough(bucket string, object string, w io.Writer) error {
	info, err := s.GrantDownloadInfo(bucket, object, 60)
	if err != nil {
		return convertErr(err)
	}
	_, err = s.PresetSave.Concat(fetcher.Get, fetcher.URL(info.URL)).FetchAndParse(fetcher.Should200(fetcher.AsDownload(w)))
	return convertErr(err)
}
func (s *Buckets) content(bucket string, object string, w io.Writer) error {
	_, err := s.PresetContent.Concat(fetcher.PathJoin(bucket), fetcher.PathJoin(object)).FetchAndParse(fetcher.Should200(fetcher.AsDownload(w)))
	if err != nil {
		return convertErr(err)
	}
	return nil
}
func (s *Buckets) Content(bucket string, object string, w io.Writer) error {
	if s.Passthrough {
		return s.contentPassthrough(bucket, object, w)
	}
	return s.content(bucket, object, w)
}
func (s *Buckets) saveFeteh(m fetcher.Method, info *WebuploadInfo, body io.Reader) (*SaveResult, error) {
	_, err := s.PresetSave.Concat(m, fetcher.URL(info.UploadURL), fetcher.Body(body)).FetchAndParse(fetcher.ShouldBetween(info.SuccessCodeMin, info.SuccessCodeMax, fetcher.AsUselessBody))
	if err != nil {
		return nil, convertErr(err)
	}
	var result = &SaveResult{
		ID:     info.ID,
		Bucket: info.Bucket,
		Object: info.Object,
	}
	return result, nil
}
func (s *Buckets) savePassthrough(bucket string, filename string, body LenReader) (*SaveResult, error) {
	opt := &UploadOptions{
		Filename:  filename,
		TTL:       strconv.Itoa(60),
		SizeLimit: strconv.Itoa(body.Len()),
	}
	info, err := s.GrantUploadInfo(bucket, opt)
	if err != nil {
		return nil, convertErr(err)
	}
	switch info.UploadType {
	case UploadTypePost:
		return s.saveFeteh(fetcher.Post, info, body)
	case UploadTypePut:
		return s.saveFeteh(fetcher.Put, info, body)
	}
	return nil, ErrUnknownUploadType
}
func (s *Buckets) save(bucket string, filename string, body LenReader) (*SaveResult, error) {
	var result = &SaveResult{}
	var writer = fetcher.NewMultiPartWriter()
	err := writer.WriteFile(PostFieldFile, filename, body)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add(QueryFieldFilename, filename)
	_, err = s.PresetSave.Concat(fetcher.PathJoin(bucket), writer, fetcher.Params(params)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, convertErr(err)
	}
	return result, nil

}
func (s *Buckets) Save(bucket string, filename string, body LenReader) (*SaveResult, error) {
	if s.Passthrough {
		return s.savePassthrough(bucket, filename, body)
	}
	return s.save(bucket, filename, body)
}
func (s *Buckets) Info(bucket string, object string) (*Fileinfo, error) {
	var result = &Fileinfo{}
	_, err := s.PresetInfo.Concat(fetcher.PathSuffix(object)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, convertErr(err)
	}
	return result, nil

}
