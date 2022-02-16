package herbbucketsgo

import (
	"bytes"
	"net/url"
	"strconv"

	"github.com/herb-go/fetcher"
)

type SaveResult struct {
	ID     string
	Bucket string
	Object string
}

type WebuploadInfo struct {
	ID             string
	UploadURL      string
	PreviewURL     string
	Permanent      bool
	Bucket         string
	Object         string
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

type Buckets struct {
	PresetGrantUploadInfo  *fetcher.Preset
	PresetGrantDownladInfo *fetcher.Preset
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

func (s *Buckets) Content(bucket string, object string) ([]byte, error) {
	var data = []byte{}
	_, err := s.PresetContent.Concat(fetcher.PathJoin(bucket), fetcher.PathJoin(object)).FetchAndParse(fetcher.Should200(fetcher.AsBytes(&data)))
	if err != nil {
		return nil, convertErr(err)
	}
	return data, nil
}

func (s *Buckets) Save(bucket string, filename string, data []byte) (*SaveResult, error) {
	var result = &SaveResult{}
	var writer = fetcher.NewMultiPartWriter()
	err := writer.WriteFile(PostFieldFile, filename, bytes.NewBuffer(data))
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

func (s *Buckets) Info(bucket string, object string) (*Fileinfo, error) {
	var result = &Fileinfo{}
	_, err := s.PresetInfo.Concat(fetcher.PathSuffix(object)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, convertErr(err)
	}
	return result, nil

}
