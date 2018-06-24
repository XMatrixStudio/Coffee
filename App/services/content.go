package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/core/errors"
)

// ContentService 内容
type ContentService interface {
	AddText(ownID, title, text string, isPublic bool, tags []string) error
	GetContentByOwn(ownID string) []models.Content
	DeleteContentByID(id, userID string) error
	PatchContentByID(id, title, content string, tags []string, public bool) error
	AddCommentCount(id string, num int) error
	AddLikeCount(id string, num int) error
}

type contentService struct {
	Model   *models.ContentModel
	Service *Service
}

func (c *contentService) AddText(ownID, title, text string, isPublic bool, tags []string) (err error) {
	_, err = c.Model.AddContent(models.Content{
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
		c.Service.Tag.AddTag(ownID, tags[i])
	}
	return
}

func (c *contentService) GetContentByOwn(ownID string) []models.Content {
	return c.Model.GetContentByOwn(ownID)

}

func (c *contentService) DeleteContentByID(id, userID string) error {
	content, err := c.Model.GetContentByID(id)
	if err != nil {
		return err
	}
	if content.OwnID.Hex() != userID {
		return errors.New("not_allow")
	}
	err = c.Model.DeleteByID(id)
	if err != nil {
		return err
	}
	// 删除评论
	c.Service.Comment.Model.DeleteAllByContent(id)
	return nil
}

func (c *contentService) PatchContentByID(id, title, content string, tags []string, public bool) error {
	con, err := c.Model.GetContentByID(id)
	if err != nil {
		return err
	}
	con.Name = title
	con.Detail = content
	con.Tag = tags
	con.Public = public
	return c.Model.UpdateByID(id, con)
}

func (c *contentService) AddCommentCount(id string, num int) error {
	return c.Model.AddCommentCount(id, num)
}

func (c *contentService) AddLikeCount(id string, num int) error {
	return c.Model.AddLikeCount(id, num)
}
