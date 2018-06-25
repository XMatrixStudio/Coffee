package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/core/errors"
)

// ContentService 内容
type ContentService interface {
	AddText(ownID, title, text string, isPublic bool, tags []string) error
	GetTextByUser(ownID string, public bool) []models.Content

	GetContentByOwn(ownID string) []models.Content
	GetContentByID(id string) (models.Content, error)
	GetPublic(int, int) []PublishData
	DeleteContentByID(id, userID string) error
	PatchContentByID(id, title, content string, tags []string, public bool) error

	AddCommentCount(id string, num int) error
	AddLikeCount(id string, num int) error

	GetUserBaseInfo(id string) (user UserBaseInfo)
}

type contentService struct {
	Model   *models.ContentModel
	Service *Service
}

func (s *contentService) AddText(ownID, title, text string, isPublic bool, tags []string) (err error) {
	_, err = s.Model.AddContent(models.Content{
		OwnID:  bson.ObjectIdHex(ownID),
		Name:   title,
		Detail: text,
		Public: isPublic,
		Tag:    tags,
		Type:   models.TypeText,
	})
	if err != nil {
		return
	}
	for i := range tags {
		s.Service.Tag.AddTag(ownID, tags[i])
	}
	return
}

func (s *contentService) GetContentByID(id string) (models.Content, error) {
	return s.Model.GetContentByID(id)
}

func (s *contentService) GetTextByUser(ownID string, public bool) []models.Content {
	return s.Model.GetContentByOwnAndType(ownID, models.TypeText)
}

func (s *contentService) GetContentByOwn(ownID string) []models.Content {
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

func (s *contentService) GetUserBaseInfo(id string) (user UserBaseInfo) {
	return s.Service.User.GetUserBaseInfo(id)
}


type PublishData struct {
	Data models.Content
	User UserBaseInfo
}

func (s *contentService) GetPublic(page, pageSize int) (contents []PublishData) {
	content := s.Model.GetPageContent(page, pageSize)
	for i := range content {
		contents = append(contents, PublishData{
			Data: content[i],
			User: s.GetUserBaseInfo(content[i].OwnID.Hex()),
		})
	}
	return
}