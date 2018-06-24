package services

import "github.com/XMatrixStudio/Coffee/App/models"

// NotificationService 通知服务
type NotificationService interface {
	AddNotification(id, content, sID, tID, typeN string) error
	GetUserNotification(userID string) (messages []NotificationData, err error)
	GetUserUnreadCount(userID string) (count int, err error)
	SetRead(userId, id string ,isRead bool) error
	RemoveByID(userID, id string) error
}

type notificationService struct {
	Model   *models.NotificationModel
	Service *Service
}

func (s *notificationService) AddNotification(userId, content, sID, tID, typeN string) error {
	return s.Model.AddNotification(content, userId, sID, tID, typeN)
}

type NotificationData struct {
	User UserBaseInfo
	Data models.NotificationDetail
}

func (s *notificationService) GetUserNotification(userID string) (data []NotificationData, err error) {
	messages, err := s.Model.GetNotificationsByUser(userID)
	for i := range messages {
		data = append(data, NotificationData{
			User: s.Service.User.GetUserBaseInfo(messages[i].SourceID),
			Data: messages[i],
		})
	}
	return
}

func (s *notificationService) GetUserUnreadCount(userID string) (count int, err error) {
	count, err = s.Model.GetUnreadCountByUser(userID)
	return
}

func (s *notificationService) SetRead(userID, id string ,isRead bool) error {
	return s.Model.ReadANotification(userID, id, isRead)
}

func (s *notificationService) RemoveByID(userID, id string) error {
	return s.Model.RemoveANotification(userID, id)
}