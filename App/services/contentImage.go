package services

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
)

func (s *contentService) GetAlbumByUser(ownID string, public bool) []models.Content {
	return s.Model.GetContentByOwnAndType(ownID, models.TypeAlbum, public)
}

func (s *contentService) AddAlbum(ctx iris.Context, id string, data ContentData) error {
	// 生成资源文件夹
	dirPath, err := s.getPath(id, models.TypeAlbum)
	if err != nil {
		return err
	}
	// 生成临时文件夹
	tmpDir := dirPath + "/" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(100000))
	if err := pathExistsAndCreate(tmpDir); err != nil {
		return err
	}
	// 生成缩略图文件夹
	thumbDir := s.getThumbDir()
	if err := pathExistsAndCreate(thumbDir); err != nil {
		return err
	}

	// 保存文件到临时文件夹
	fileSize, err := ctx.UploadFormFiles(tmpDir, s.dealWithAlbum)
	defer os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}
	// 增加用户已用存储
	if err := s.Service.User.AddFiles(id, fileSize); err != nil {
		// 用户存储不足
		return err
	}
	// 处理文件
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	var Images []models.Image
	for _, f := range files {
		name := strings.Split(f.Name(), "-")
		if len(name) < 2 {
			return errors.New("error_name")
		}
		if strings.Contains(name[1], "file") {
			index := strings.Split(name[1], "file")
			if len(index) != 2 {
				return errors.New("error_name")
			}
			thumbFileIndex := indexOfFile(files, "thumb"+index[1])
			if thumbFileIndex == -1 {
				return errors.New("error_name")
			}
			fileTargetPath := dirPath + "/" + f.Name()
			if _, err := copyFile(tmpDir+"/"+f.Name(), fileTargetPath); err != nil {
				return err
			}
			if _, err := copyFile(tmpDir+"/"+files[thumbFileIndex].Name(), thumbDir+"/"+name[0]+".png"); err != nil {
				return err
			}
			Images = append(Images, models.Image{
				Native: true,
				File: models.File{
					File:  fileTargetPath,
					Size:  f.Size(),
					Title: strings.Replace(f.Name(), name[0]+"-"+name[1]+"-", "", -1),
					Time:  time.Now().Unix() * 1000,
					Type:  models.TypeAlbum,
				},
				Thumb: name[0] + ".png",
			})
		}
	}

	_, err = s.Model.AddContent(models.Content{
		OwnID:  bson.ObjectIdHex(id),
		Name:   data.Title,
		Detail: data.Detail,
		Public: data.IsPublic,
		Tag:    data.Tags,
		Type:   models.TypeAlbum,
		Album: models.Album{
			Images: Images,
			Title:  data.Title,
			Time:   time.Now().Unix() * 1000,
		},
	})

	return err
}

// dealWithFile 处理文件
func (s *contentService) dealWithAlbum(ctx iris.Context, file *multipart.FileHeader) {
	re, _ := regexp.Compile(`form-data; name="(.*)"; filename=`)
	var typeName string
	if match := re.FindStringSubmatch(file.Header.Get("Content-Disposition")); len(match) != 2 {
		typeName = "Unknown"
	} else {
		typeName = match[1]
	}
	if strings.Contains(typeName, "thumb") {
		typeName += ".png"
	} else {
		typeName += "-" + file.Filename
	}
	file.Filename = strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(100000)) + "-" + typeName
}

func (s *contentService) getPath(userID, fileType string) (string, error) {
	dirPath := s.UserDir + "/" + userID + "/" + fileType
	if err := pathExistsAndCreate(dirPath); err != nil {
		return "", err
	}
	return dirPath, nil
}

func (s *contentService) getThumbDir() string {
	dirPath := s.ThumbDir
	pathExistsAndCreate(dirPath)
	return dirPath
}
