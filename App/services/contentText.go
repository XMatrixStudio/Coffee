package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/globalsign/mgo/bson"
)

func (s *contentService) AddText(ownID string, data ContentData) (err error) {
	_, err = s.Model.AddContent(models.Content{
		OwnID:  bson.ObjectIdHex(ownID),
		Name:   data.Title,
		Detail: data.Detail,
		Public: data.IsPublic,
		Tag:    data.Tags,
		Type:   models.TypeText,
	})
	if err != nil {
		return
	}
	for i := range data.Tags {
		s.Service.Tag.AddTag(ownID, data.Tags[i])
	}
	return
}

func (s *contentService) GetTextByUser(ownID string, public bool) []models.Content {
	return s.Model.GetContentByOwnAndType(ownID, models.TypeText, public)
}