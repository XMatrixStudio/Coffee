package services

import "github.com/XMatrixStudio/Coffee/App/models"

// FollowService 关注
type FollowService interface {
	GetUserFollower(userID string) (ids []string, err error)
	GetUserFollowing(userID string) (ids []string, err error)
	AddFollow(userFollowing, userFollower string) error
	RemoveFollow(userFollowing, userFollower string) error
}

type followService struct {
	Model   *models.GatherModel
	Service *Service
}

func (s *followService) GetUserFollower(userID string) (ids []string, err error) {
	return s.Model.GetFollower(userID)
}

func (s *followService) GetUserFollowing(userID string) (ids []string, err error) {
	return s.Model.GetFollowing(userID)
}

func (s *followService) AddFollow(userID, followID string) error {
	if err := s.Model.AddUserFollowing(userID, followID); err != nil {
		return err
	}
	if err := s.Model.AddUserFollower(followID, userID); err != nil {
		s.Model.RemoveFollowing(userID, followID)
		return err
	}
	// todo
	return nil
}

func (s *followService) RemoveFollow(userFollowing, userFollower string) error{
	// todo
	return nil
}
