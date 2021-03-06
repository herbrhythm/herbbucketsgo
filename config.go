package herbbucketsgo

import (
	"net/http"
	"net/url"
	"path"

	"github.com/herb-go/fetcher"
)

func join(u url.URL, ele string) string {
	u.Path = path.Join(u.Path, ele)
	return u.String()
}

type Config struct {
	*fetcher.Server
	Passthrough          bool
	Bucket               string
	URLGrantUploadInfo   string
	URLGrantDownloadInfo string
	URLComplete          string
	URLContent           string
	URLSave              string
	URLRemove            string
	URLInfo              string
}

func (c Config) Apply(b *Buckets) error {
	var err error
	u, err := url.Parse(c.Server.URL)
	if err != nil {
		return err
	}
	if c.URLGrantUploadInfo == "" {
		c.URLGrantUploadInfo = join(*u, "/api/grantuploadinfo")
	}
	b.PresetGrantUploadInfo, err = c.Server.MergeURL(c.URLGrantUploadInfo).MergeMethod(http.MethodPost).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLGrantDownloadInfo == "" {
		c.URLGrantDownloadInfo = join(*u, "/api/grantdownloadinfo")
	}
	b.PresetGrantDownladInfo, err = c.Server.MergeURL(c.URLGrantDownloadInfo).MergeMethod(http.MethodPost).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLComplete == "" {
		c.URLComplete = join(*u, "/complete")
	}
	b.PresetComplete, err = c.Server.MergeURL(c.URLComplete).MergeMethod(http.MethodPost).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLContent == "" {
		c.URLContent = join(*u, "/api/content")
	}
	b.PresetContent, err = c.Server.MergeURL(c.URLContent).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLSave == "" {
		c.URLSave = c.Server.URL + "/api/content"
	}
	b.PresetSave, err = c.Server.MergeURL(c.URLSave).MergeMethod(http.MethodPost).CreatePreset()
	if err != nil {
		return err
	}

	if c.URLRemove == "" {
		c.URLRemove = join(*u, "/api/remove")
	}
	b.PresetRemove, err = c.Server.MergeURL(c.URLRemove).MergeMethod(http.MethodPost).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLInfo == "" {
		c.URLInfo = join(*u, "/api/remove")
	}
	b.PresetInfo, err = c.Server.MergeURL(c.URLInfo).MergeMethod(http.MethodPost).CreatePreset()
	if err != nil {
		return err
	}
	b.Bucket = c.Bucket
	b.Passthrough = c.Passthrough
	return nil
}
