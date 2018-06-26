package services

import (
	"github.com/kataras/iris"
	"mime/multipart"
	"fmt"
	"os"
	"time"
	"strconv"
	"math/rand"
	"regexp"
	"strings"
	"github.com/XMatrixStudio/Coffee/App/models"
	"io/ioutil"
	"github.com/kataras/iris/core/errors"
	"io"
)

func (s *contentService) AddAlbum(ctx iris.Context, id string) error {
	dirPath, err := getPath(id, models.TypeAlbum)
	if err != nil {
		return err
	}
	// 生成临时文件夹
	tmpDir := dirPath + "/" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(rand.Intn(100000))
	if err := pathExistsAndCreate(tmpDir); err != nil {
		return  err
	}
	// 保存文件到临时文件夹
	fileSize, err := ctx.UploadFormFiles(tmpDir, dealWithFile)
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
	for _, f := range files {
		fmt.Println(f.Name())
		name := strings.Split(f.Name(), "-")
		if len(name) < 2 {
			return errors.New("error_name")
		}

		if strings.Contains(name[1], "thumb") {
			if _, err := copyFile(tmpDir + "/" + f.Name(), getThumbDir() + "/" + f.Name()); err != nil {
				return err
			}
		} else {
			if _, err := copyFile(tmpDir + "/" + f.Name(), dirPath + "/" + f.Name()); err != nil {
				return err
			}
		}
	}

	return err
}

func copyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// dealWithFile 处理文件
func dealWithFile(ctx iris.Context, file *multipart.FileHeader) {
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
	fmt.Println(file.Filename)
}

func getPath(userID, fileType string) (string, error) {
	dirPath := "userData" +  "/" + userID + "/" + fileType
	if err := pathExistsAndCreate(dirPath); err != nil {
		return "", err
	}
	return dirPath, nil
}

func getThumbDir() string {
	dirPath := "thumb"
	if err := pathExistsAndCreate(dirPath); err != nil {
		fmt.Errorf(err.Error())
	}
	return dirPath
}

// 判断文件夹是否存在
func pathExistsAndCreate(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}