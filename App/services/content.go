package services

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
)

// ContentService 内容
type ContentService interface {
	AddText(ownID, title, text string, isPublic bool, tags []string) error
	GetTextByUser(ownID string, public bool) []models.Content

	GetContentsByOwn(ownID string) []models.Content
	GetContentByID(id string) (models.Content, error)
	GetContentAndUser(id string) (content models.Content, user UserBaseInfo, err error)
	GetPublicContents(int, int) []PublishData
	DeleteContentByID(id, userID string) error
	PatchContentByID(id, title, content string, tags []string, public bool) error

	AddCommentCount(id string, num int) error
	AddLikeCount(id string, num int) error

	BeforeSave(ctx iris.Context, file *multipart.FileHeader)
}

type contentService struct {
	Model   *models.ContentModel
	Service *Service
}

func (s *contentService) GetContentByID(id string) (models.Content, error) {
	return s.Model.GetContentByID(id)
}

func (s *contentService) GetContentsByOwn(ownID string) []models.Content {
	return s.Model.GetContentByOwn(ownID)
}

func (s *contentService) DeleteContentByID(id, userID string) error {
	content, err := s.Model.GetContentByID(id)
	if err != nil {
		return err
	}
	if content.OwnID.Hex() != userID {
		return errors.New("not_allow")
	}
	err = s.Model.DeleteByID(id)
	if err != nil {
		return err
	}
	// 删除评论
	s.Service.Comment.Model.DeleteAllByContent(id)
	s.Service.Like.Model.RemoveAllByID(id)
	return nil
}

func (s *contentService) PatchContentByID(id, title, content string, tags []string, public bool) error {
	con, err := s.Model.GetContentByID(id)
	if err != nil {
		return err
	}
	con.Name = title
	con.Detail = content
	con.Tag = tags
	con.Public = public
	return s.Model.UpdateByID(id, con)
}

func (s *contentService) AddCommentCount(id string, num int) error {
	return s.Model.AddCommentCount(id, num)
}

func (s *contentService) AddLikeCount(id string, num int) error {
	return s.Model.AddLikeCount(id, num)
}

func (s *contentService) GetContentAndUser(id string) (content models.Content, user UserBaseInfo, err error) {
	content, err = s.GetContentByID(id)
	if err != nil {
		return
	}
	user = s.Service.User.GetUserBaseInfo(content.OwnID.Hex())
	return
}

// PublishData 公共数据
type PublishData struct {
	Data models.Content
	User UserBaseInfo
}

// GetPublicContents 获取公共内容
func (s *contentService) GetPublicContents(page, pageSize int) (contents []PublishData) {
	content := s.Model.GetPageContent(page, pageSize)
	for i := range content {
		contents = append(contents, PublishData{
			Data: content[i],
			User: s.Service.User.GetUserBaseInfo(content[i].OwnID.Hex()),
		})
	}
	return
}

// BeforeSave 处理文件
func (s *contentService) BeforeSave(ctx iris.Context, file *multipart.FileHeader) {

	ip := ctx.RemoteAddr()
	// make sure you format the ip in a way
	// that can be used for a file name (simple case):
	ip = strings.Replace(ip, ".", "_", -1)

	// you can use the time.Now, to prefix or suffix the files
	// based on the current time as well, as an exercise.
	// i.e unixTime :=	time.Now().Unix()
	// prefix the Filename with the $IP-
	// no need for more actions, internal uploader will use this
	// name to save the file into the "./uploads" folder.
	file.Filename = ip + "-" + file.Filename

	fmt.Println(file.Size)
}
