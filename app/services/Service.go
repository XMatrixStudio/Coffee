package services

import "github.com/XMatrixStudio/Coffee/App/models"

// Service 服务
type Service struct {
	Model        *models.Model
	User         userService
	Notification notificationService
	Tag          tagService
	Content      contentService
	Comment      commentService
	Like likeService
}

// NewService 初始化Service
func NewService(m *models.Model) *Service {
	service := new(Service)
	service.Model = m
	service.User = userService{
		Model:    &m.User,
		Service:  service,
		UserInfo: make(map[string]UserBaseInfo),
	}
	service.Notification = notificationService{
		Model:   &m.Notification,
		Service: service,
	}
	service.Tag = tagService{
		Model:   &m.Tag,
		Service: service,
	}
	service.Content = contentService{
		Model:   &m.Content,
		Service: service,
	}
	service.Comment = commentService{
		Model: &m.Comment,
		Service: service,
	}
	service.Like = likeService{
		Model: &m.Gather,
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

// GetContentService 新建 ContentService
func (s *Service) GetContentService() ContentService {
	return &s.Content
}

// GetCommentService 评论
func (s *Service) GetCommentService() CommentService {
	return &s.Comment
}

// GetLikeService 点赞
func (s *Service) GetLikeService() LikeService {
	return &s.Like
}