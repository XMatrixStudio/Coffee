package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/XMatrixStudio/Coffee/App/models"
	"io"
	"mime/multipart"
	"os"
)

type FileService interface {
	InitFileService(tempPath, dataPath string)
	AddOwn(token, ownID string) error
	AddFile(data UploadMeta) error
	GetInfo(md5 string) (res UploadMeta, err error)
	UploadFileToTemp(file *multipart.File, token string, index int) error
	MergeFile(token string, userID string) error
	IsExist(token string) bool
}

type fileService struct {
	Service   *Service
	Model     *models.FileModel
	Uploading map[string]UploadMeta
	TempDir   string
	DataDir   string
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

func (s *fileService) InitFileService(tempPath, dataPath string) {
	s.TempDir = tempPath
	pathExistsAndCreate(tempPath)
	s.DataDir = dataPath
	pathExistsAndCreate(dataPath)
	s.Uploading = make(map[string]UploadMeta)
}

func (s *fileService) AddFile(data UploadMeta) error {
	for i := range data.Chunks {
		data.Chunks[i].Success = false
	}
	s.Uploading[data.MD5] = data
	fmt.Println(data)
	return nil
}

func (s *fileService) AddOwn(token, ownID string) error {
	return s.Model.AddOwn(ownID, token)
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

func (s *fileService) IsExist(token string) bool {
	return s.Model.IsExist(token)
}

func (s *fileService) UploadFileToTemp(file *multipart.File, token string, index int) error {
	if meta, ok := s.Uploading[token]; !ok {
		return errors.New("not_exist")
	} else if index < 0 || index >= len(meta.Chunks) {
		return errors.New("not_exist")
	}

	// 检测合法性
	hash := md5.New()
	if _, err := io.Copy(hash, *file); err != nil {
		return err
	}
	hashInBytes := hash.Sum(nil)[:16]
	MD5String := hex.EncodeToString(hashInBytes)
	fmt.Println(MD5String, s.Uploading[token].Chunks[index].MD5)
	if s.Uploading[token].Chunks[index].MD5 != MD5String {
		return errors.New("invalid_chunk")
	}

	out, err := os.OpenFile(s.TempDir+"/"+MD5String, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer out.Close()
	(*file).Seek(0, 0)
	if _, err := io.Copy(out, *file); err != nil {
		return err
	}

	s.Uploading[token].Chunks[index].Success = true
	return nil
}

func (s *fileService) MergeFile(token string, userID string) error {
	if meta, ok := s.Uploading[token]; !ok {
		return errors.New("not_exist")
	} else {
		for _, chunk := range meta.Chunks {
			if chunk.Success != true {
				return errors.New("lack_chunk")
			}
		}
		path := s.DataDir + "/" + userID
		pathExistsAndCreate(path)
		out, err := os.OpenFile(path+"/"+meta.MD5, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer out.Close()
		for _, chunk := range meta.Chunks {
			if data, err := os.OpenFile(s.TempDir+"/"+chunk.MD5, os.O_RDONLY, os.ModePerm); err != nil {
				return err
			} else {
				io.Copy(out, data)
				data.Close()
			}
		}
	}
	delete(s.Uploading, token)
	return nil
}
