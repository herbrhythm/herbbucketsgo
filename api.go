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
	Sizelimit      int64
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
type Buckets struct {
	PresetGrantUploadInfo  *fetcher.Preset
	PresetGrantDownladInfo *fetcher.Preset
	PresetContent          *fetcher.Preset
	PresetSave             *fetcher.Preset
	PresetRemove           *fetcher.Preset
	PresetInfo             *fetcher.Preset
}

func (s *Buckets) GrantUploadInfo(opt *UploadOptions) (*WebuploadInfo, error) {
	result := &WebuploadInfo{}
	params := url.Values{}
	params.Add(QueryFieldFilename, opt.Filename)
	params.Add(QueryFieldTTL, opt.TTL)
	params.Add(QueryFieldSizelimit, opt.Sizelimit)
	params.Add(QueryFieldSize, opt.Size)
	_, err := s.PresetGrantUploadInfo.Concat(fetcher.Params(params)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Buckets) GrantDownloadInfo(object string, ttl int64) (*DownloadInfo, error) {
	var result = &DownloadInfo{}
	params := url.Values{}
	params.Add(QueryFieldTTL, strconv.FormatInt(ttl, 10))
	_, err := s.PresetGrantDownladInfo.Concat(fetcher.PathSuffix(object), fetcher.Params(params)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Buckets) Content(object string) ([]byte, error) {
	var data = []byte{}
	_, err := s.PresetContent.Concat(fetcher.PathSuffix(object)).FetchAndParse(fetcher.Should200(fetcher.AsBytes(&data)))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Buckets) Save(filename string, data []byte) (*SaveResult, error) {
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
	_, err = s.PresetSave.Concat(writer, fetcher.Params(params)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (s *Buckets) Info(object string) (*Fileinfo, error) {
	var result = &Fileinfo{}
	_, err := s.PresetInfo.Concat(fetcher.PathSuffix(object)).FetchAndParse(fetcher.Should200(fetcher.AsJSON(result)))
	if err != nil {
		return nil, err
	}
	return result, nil

}
