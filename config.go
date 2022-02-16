package herbbucketsgo

import (
	"path"

	"github.com/herb-go/fetcher"
)

type Config struct {
	*fetcher.Server
	URLGrantUploadInfo   string
	URLGrantDownloadInfo string
	URLContent           string
	URLSave              string
	URLRemove            string
	URLInfo              string
}

func (c Config) Apply(b *Buckets) error {
	var err error
	if c.URLGrantUploadInfo == "" {
		c.URLGrantUploadInfo = path.Join(c.Server.URL, "api", "grantuploadinfo")
	}
	b.PresetGrantUploadInfo, err = c.Server.MergeURL(c.URLGrantUploadInfo).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLGrantDownloadInfo == "" {
		c.URLGrantDownloadInfo = path.Join(c.Server.URL, "api", "grantdownloadinfo")
	}
	b.PresetGrantDownladInfo, err = c.Server.MergeURL(c.URLGrantDownloadInfo).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLContent == "" {
		c.URLContent = path.Join(c.Server.URL, "api", "content")
	}
	b.PresetContent, err = c.Server.MergeURL(c.URLContent).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLSave == "" {
		c.URLSave = path.Join(c.Server.URL, "api", "content")
	}
	b.PresetSave, err = c.Server.MergeURL(c.URLSave).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLRemove == "" {
		c.URLRemove = path.Join(c.Server.URL, "api", "remove")
	}
	b.PresetRemove, err = c.Server.MergeURL(c.URLRemove).CreatePreset()
	if err != nil {
		return err
	}
	if c.URLInfo == "" {
		c.URLInfo = path.Join(c.Server.URL, "api", "remove")
	}
	b.PresetInfo, err = c.Server.MergeURL(c.URLInfo).CreatePreset()
	if err != nil {
		return err
	}
	return nil
}
