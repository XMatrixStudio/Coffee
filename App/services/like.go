package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
)

// LikeService 点赞
type LikeService interface {
	GetUserLikes(userID string) (ids []string, err error)
	AddLikeToContent(id, userID string) error
	AddLikeToComment(id, userID string, isReply bool) error
	RemoveLikeFromContent(id, userID string) error
	RemoveLikeFromComment(id, userID string, isReply bool) error
}

type likeService struct {
	Model   *models.GatherModel
	Service *Service
}

func (s *likeService) AddLikeToContent(id, userID string) error {
	err := s.Model.AddLikeToUser(id, userID)
	if err != nil {
		// 已经点赞了
		return err
	}
	err = s.Service.Content.AddLikeCount(id, 1)
	if err != nil {
		// 增加失败，需要回退结果
		s.Model.RemoveLikeFromUser(id, userID)
		return err
	} else {
		// 冗余数据
		s.Model.AddLikeToContent(id, userID)
		// 给作者发送通知
		content, err := s.Service.Content.Model.GetContentByID(id)
		if err == nil && content.OwnID.Hex() != userID {
			s.Service.Notification.AddNotification(content.OwnID.Hex(), "", userID, id, models.TypeLike)
		}
	}
	return nil
}

func (s *likeService) RemoveLikeFromContent(id, userID string) error {
	err := s.Model.RemoveLikeFromUser(id, userID)
	if err != nil {
		// 不存在点赞内容
		return err
	}
	s.Service.Content.AddLikeCount(id, -1)
	s.Model.RemoveLikeFromContent(id, userID)
	// 删除给作者发送的未读通知
	content, err := s.Service.Content.Model.GetContentByID(id)
	if err == nil {
		s.Service.Notification.Model.RemoveUnread(content.OwnID.Hex(), userID, id)
	}
	return nil
}

func (s *likeService) AddLikeToComment(id, userID string, isReply bool) error {
	err := s.Model.AddLikeToUser(id, userID)
	if err != nil {
		// 已经点赞了
		return err
	}
	if isReply {
		err = s.Service.Comment.Model.AddLikeToReply(id, 1)
	} else {
		err = s.Service.Comment.Model.AddLikeToComment(id, 1)
	}
	if err != nil {
		// 增加失败，需要回退结果
		s.Model.RemoveLikeFromUser(id, userID)
	}
	return err
}

func (s *likeService) RemoveLikeFromComment(id, userID string, isReply bool) error {
	err := s.Model.RemoveLikeFromUser(id, userID)
	if err != nil {
		// 不存在点赞内容
		return err
	}
	if isReply {
		s.Service.Comment.Model.AddLikeToReply(id, -1)
	} else {
		s.Service.Comment.Model.AddLikeToComment(id, -1)
	}
	return nil
}

func (s *likeService) GetUserLikes(userID string) (ids []string, err error) {
	ids, err = s.Model.GetUserLikes(userID)
	return
}