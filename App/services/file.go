package services

import "github.com/kataras/iris/core/errors"

type FileService interface {
	AddFile(data UploadMeta) error
}

type fileService struct {
	Service   *Service
	Uploading map[string]UploadMeta
}

type UploadMeta struct {
	Name   string      `json:"name"`
	Size   int64       `json:"size"`
	MD5    string      `json:"md5"`
	Chunks []FileChunk `json:"chunks"`
}

type FileChunk struct {
	Index   int64  `json:"index"`
	Size    int64  `json:"size"`
	MD5     string `json:"md5"`
	Success bool   `json:"success"`
}

func (s *fileService) AddFile(data UploadMeta) error {
	if _, ok := s.Uploading[data.MD5]; ok {
		return errors.New("exist")
	}
	s.Uploading[data.MD5] = data
	return nil
}

func (s *fileService) GetInfo(md5 string) (res UploadMeta, err error) {
	if meta, ok := s.Uploading[md5]; ok {
		res = meta
		return
	} else {
		err = errors.New("not_exist")
		return
	}
}
