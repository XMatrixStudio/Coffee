package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/globalsign/mgo/bson"
)

// ContentService 内容
type ContentService interface {
	AddText(ownID, title, text string, isPublic bool, tags []string) (err error)
	GetContentByOwn(ownID string) (data []models.Content)
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
	})
	if err != nil {
		return
	}
	for i := range tags {
		c.Service.Tag.AddTag(ownID, tags[i])
	}
	return
}

func (c *contentService) GetContentByOwn(ownID string) (data []models.Content) {
	data = c.Model.GetContentByOwn(ownID)
	return
}
