package services

import "github.com/XMatrixStudio/Coffee/App/models"

// TagService 通知服务
type TagService interface {
	AddNotification(id, tag string) error
}

type tagService struct {
	Model   *models.TagModel
	Service *Service
}

func (s *tagService) AddTag(id, tag string) error {
	return s.Model.AddTag(id, tag)
}
