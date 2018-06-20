package services

import "github.com/XMatrixStudio/Coffee/App/models"

// NotificationService 通知服务
type NotificationService interface {
	AddNotification(id, content, typeN string) error
}

type notificationService struct {
	Model   *models.NotificationModel
	Service *Service
}

func (s *notificationService) AddNotification(id, content, typeN string) error {
	return s.Model.AddNotification(content, id, typeN)
}
