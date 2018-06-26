package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/globalsign/mgo/bson"
)

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

func (s *contentService) GetTextByUser(ownID string, public bool) []models.Content {
	return s.Model.GetContentByOwnAndType(ownID, models.TypeText, public)
}