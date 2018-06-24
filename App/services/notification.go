package services

import "github.com/XMatrixStudio/Coffee/App/models"

// NotificationService 通知服务
type NotificationService interface {
	AddNotification(id, content, typeN string) error
	GetUserNotification(userID string) (messages []models.NotificationDetail, err error)
	GetUserUnreadCount(userID string) (count int, err error)
}

type notificationService struct {
	Model   *models.NotificationModel
	Service *Service
}

func (s *notificationService) AddNotification(id, content, sID, tID, typeN string) error {
	return s.Model.AddNotification(content, id, sID, tID, typeN)
}

func (s *notificationService) GetUserNotification(userID string) (messages []models.NotificationDetail, err error) {
	messages, err = s.Model.GetNotificationsByUser(userID)
	return
}

func (s *notificationService) GetUserUnreadCount(userID string) (count int, err error) {
	count, err = s.Model.GetUnreadCountByUser(userID)
	return
}
