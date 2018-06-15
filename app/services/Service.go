package services

import "github.com/XMatrixStudio/Coffee/app/models"

type Service struct {
	Model *models.Model
}

// NewService 初始化Service
func NewService(m *models.Model) *Service {
	service := new(Service)
	service.Model = m
	return service
}

func (s *Service) NewUserService() UserService {
	return &userService{
		Model: s.Model.User,
	}
}
