package services

import "github.com/XMatrixStudio/Coffee/App/models"

// Service 服务
type Service struct {
	Model        *models.Model
	User         userService
	Notification notificationService
	Tag          tagService
}

// NewService 初始化Service
func NewService(m *models.Model) *Service {
	service := new(Service)
	service.Model = m
	service.User = userService{
		Model:   &m.User,
		Service: service,
	}
	service.Notification = notificationService{
		Model:   &m.Notification,
		Service: service,
	}
	service.Tag = tagService{
		Model:   &m.Tag,
		Service: service,
	}
	return service
}

// GetUserService 新建 UserService
func (s *Service) GetUserService() UserService {
	return &s.User
}

// GetNotificationService 新建 NotificationService
func (s *Service) GetNotificationService() NotificationService {
	return &s.Notification
}
