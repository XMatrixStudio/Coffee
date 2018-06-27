package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/kataras/iris"
	"errors"
)

// ContentService 内容
type ContentService interface {
	SetThumbDir(path string)
	SetUserDir(path string)

	AddText(ownID, title, text string, isPublic bool, tags []string) error
	GetTextByUser(ownID string, public bool) []models.Content

	GetFile(userID, contentID, filePath string) (string, error)

	GetContentsByOwn(ownID string) []models.Content
	GetContentByID(id string) (models.Content, error)
	GetContentAndUser(id string) (content models.Content, user UserBaseInfo, err error)
	GetPublicContents(int, int) []PublishData
	DeleteContentByID(id, userID string) error
	PatchContentByID(id, title, content string, tags []string, public bool) error

	AddCommentCount(id string, num int) error
	AddLikeCount(id string, num int) error

	AddAlbum(ctx iris.Context, id string, data NewAlbumData) error
	GetAlbumByUser(ownID string, public bool) []models.Content
}

type contentService struct {
	Model   *models.ContentModel
	Service *Service
	ThumbDir string
	UserDir string
}

func (s *contentService) SetThumbDir(path string) {
	s.ThumbDir = path
}

func (s *contentService) SetUserDir (path string) {
	s.UserDir = path
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


func (s *contentService) GetFile(userID, contentID, filePath string) (string, error) {
	content, err := s.Model.GetContentByID(contentID)
	if err != nil {
		return "", err
	}
	if content.Public == false && userID != content.OwnID.Hex() {
		return "", err
	}
	return s.Model.FindFileInContent(contentID, filePath)
}