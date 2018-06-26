package models

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Comment 评论
type Comment struct {
	ID        bson.ObjectId `bson:"_id"`
	ContentID bson.ObjectId `bson:"contentId"` // 内容ID
	FatherID  bson.ObjectId `bson:"fatherId"`  // 内容用户ID
	UserID    bson.ObjectId `bson:"userId"`    // 评论用户ID
	Date      int64         `bson:"date"`      // 发布时间
	Content   string        `bson:"content"`   // 评论内容
	LikeNum   int64         `bson:"likeNum"`   // 点赞数
}

// CommentModel 评论数据库
type CommentModel struct {
	CommentDB *mgo.Collection
	ReplyDB   *mgo.Collection
}

// AddComment 增加评论
func (m *CommentModel) AddComment(contentID, userID, content, fatherID string, isReply bool) (err error) {
	var targetDB *mgo.Collection
	if isReply {
		targetDB = m.ReplyDB
	} else {
		targetDB = m.CommentDB
	}
	return targetDB.Insert(&Comment{
		ID:        bson.NewObjectId(),
		FatherID:  bson.ObjectIdHex(fatherID),
		ContentID: bson.ObjectIdHex(contentID),
		UserID:    bson.ObjectIdHex(userID),
		Content:   content,
		Date:      time.Now().Unix() * 1000,
	})
}

// AddLikeToComment 给评论点赞
func (m *CommentModel) AddLikeToComment(id string, num int) (err error) {
	return m.CommentDB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
}

// AddLikeToReply 给回复点赞
func (m *CommentModel) AddLikeToReply(id string, num int) (err error) {
	return m.ReplyDB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
}

// RemoveComment 删除评论
func (m *CommentModel) RemoveComment(id string) (err error) {
	return m.CommentDB.RemoveId(bson.ObjectIdHex(id))
}

// RemoveReply 删除评论
func (m *CommentModel) RemoveReply(id string) (err error) {
	return m.ReplyDB.RemoveId(bson.ObjectIdHex(id))
}

// GetCommentByContentID 获取内容的评论
func (m *CommentModel) GetCommentByContentID(id string) []Comment {
	var comment []Comment
	m.CommentDB.Find(bson.M{"contentId": bson.ObjectIdHex(id)}).Sort("-date").All(&comment)
	return comment
}

func (m *CommentModel) GetCommentByID(id string) (comment Comment, err error) {
	err = m.CommentDB.FindId(bson.ObjectIdHex(id)).One(&comment)
	return
}

func (m *CommentModel) GetReplyByID(id string) (reply Comment, err error) {
	err = m.ReplyDB.FindId(bson.ObjectIdHex(id)).One(&reply)
	return
}

// GetReplyByCommentID 获取指定ID评论的回复
func (m *CommentModel) GetReplyByCommentID(id string) (reply []Comment, err error) {
	err = m.ReplyDB.Find(bson.M{"contentId": bson.ObjectIdHex(id)}).Sort("date").All(&reply)
	return
}

// GetCommentCount 获取评论数目
func (m *CommentModel) GetCommentCount(id string) (count int) {
	count, err := m.CommentDB.Find(bson.M{"contentId": bson.ObjectIdHex(id)}).Count()
	if err != nil {
		count = -1
	}
	return
}

// DeleteAllByContent
func (m *CommentModel) DeleteAllByContent(id string) {
	var comment []Comment
	err := m.CommentDB.Find(bson.M{"contentId": bson.ObjectIdHex(id)}).All(&comment)
	if err != nil {
		return
	}
	// 删除所有回复
	for i := range comment {
		m.ReplyDB.RemoveAll(bson.M{"contentId": comment[i].ID})
	}
	// 删除所有评论
	m.CommentDB.RemoveAll(bson.M{"contentId": bson.ObjectIdHex(id)})
}
