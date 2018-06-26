package services

import (
	"errors"

	"github.com/XMatrixStudio/Coffee/App/models"
)

// CommentService 评论
type CommentService interface {
	GetComment(id string) (comments []CommentForContent)
	AddComment(userID, contentID, fatherID, content string, isReply bool) error
	DeleteComment(id, userID string) error
}

type commentService struct {
	Model   *models.CommentModel
	Service *Service
}

// CommentForContent 一条内容的评论和回复
type CommentForContent struct {
	Comment models.Comment
	User    UserBaseInfo
	Replies []ReplyForComment
}

// ReplyForComment 一条评论的回复
type ReplyForComment struct {
	Reply  models.Comment
	User   UserBaseInfo
	Father UserBaseInfo
}

func (s *commentService) GetComment(id string) (comments []CommentForContent) {
	commentAll := s.Model.GetCommentByContentID(id)
	for i := range commentAll {
		replyAll, err := s.Model.GetReplyByCommentID(commentAll[i].ID.Hex())
		var replies []ReplyForComment
		if err == nil {
			for j := range replyAll {
				replies = append(replies, ReplyForComment{
					Reply:  replyAll[j],
					User:   s.Service.User.GetUserBaseInfo(replyAll[j].UserID.Hex()),
					Father: s.Service.User.GetUserBaseInfo(replyAll[j].FatherID.Hex()),
				})
			}
		}
		comments = append(comments, CommentForContent{
			Comment: commentAll[i],
			User:    s.Service.User.GetUserBaseInfo(commentAll[i].UserID.Hex()),
			Replies: replies,
		})
	}
	return
}

func (s *commentService) AddComment(userID, contentID, fatherID, content string, isReply bool) (err error) {
	err = s.Model.AddComment(contentID, userID, content, fatherID, isReply)
	if err != nil {
		return
	}
	if !isReply {
		s.Service.Content.AddCommentCount(contentID, 1)
	}
	// 给father发送通知
	if fatherID != userID {
		if !isReply {
			s.Service.Notification.AddNotification(fatherID, content, userID, contentID, models.TypeReply)
		} else {
			comment, err := s.Model.GetCommentByID(contentID)
			if err != nil {
				return err
			}
			s.Service.Notification.AddNotification(fatherID, content, userID, comment.ContentID.Hex(), models.TypeReply)
		}
	}
	return
}

func (s *commentService) DeleteComment(id, userID string) error {
	isReply := false
	comment, err := s.Model.GetCommentByID(id)
	if err != nil && err.Error() == "not found" {
		comment, err = s.Model.GetReplyByID(id)
		if err != nil {
			return err
		}
		isReply = true
	} else if err != nil {
		return err
	}
	if isReply {
		if err != nil {
			return err
		}
		commentFather, err := s.Model.GetCommentByID(comment.ContentID.Hex())
		content, err := s.Service.Content.Model.GetContentByID(commentFather.ContentID.Hex())
		if err != nil {
			return err
		}
		// 是否有删除权限(回复者和内容拥有者)
		if comment.UserID.Hex() != userID && content.OwnID.Hex() != userID {
			return errors.New("not_allow")
		}
		s.Model.RemoveReply(id)
	} else {
		if err != nil {
			return err
		}
		// 是否有删除权限(评论者和内容拥有者)
		if comment.UserID.Hex() != userID && comment.FatherID.Hex() != userID {
			return errors.New("not_allow")
		}
		// 删除评论下所有回复
		replies, err := s.Model.GetReplyByCommentID(id)
		if err == nil {
			for i := range replies {
				s.Model.RemoveReply(replies[i].ID.Hex())
			}
		}
		s.Model.RemoveComment(id)
		s.Service.Content.AddCommentCount(comment.ContentID.Hex(), -1)
	}
	s.Service.Like.Model.RemoveAllByID(id)
	return nil
}
