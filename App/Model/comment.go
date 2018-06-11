package model

import (
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Comment 评论
type Comment struct {
	ID        bson.ObjectId `bson:"_id"`
	ContentID string        `bson:"contentId"` // 文章ID 【索引】
	UserID    string        `bson:"userId"`    // 评论用户ID 【索引】
	Date      int64         `bson:"date"`      // 发布时间
	Content   string        `bson:"content"`   // 评论内容
	LikeNum   int64         `bson:"likeNum"`   // 点赞数
}

// Reply 评论的回复
type Reply struct {
	ID        bson.ObjectId `bson:"_id"`
	CommentID string        `bson:"commentId"` // 父评论ID序列
	UserID    string        `bson:"userId"`    // 评论用户ID
	FatherID  string        `bson:"fatherId"`  // 被评论用户ID
	Date      int64         `bson:"date"`      // 发布时间
	Content   string        `bson:"content"`   // 评论内容
	LikeNum   int64         `bson:"likeNum"`   // 点赞数
}

// CommentDB 评论数据库
var CommentDB *mgo.Collection

// ReplyDB 回复数据库
var ReplyDB *mgo.Collection

// AddComment 增加评论
func AddComment(contentID, userID, content, fatherID string) (bson.ObjectId, error) {
	newComment := bson.NewObjectId()
	if fatherID == "" {
		err := CommentDB.Insert(&Comment{
			ID:        newComment,
			ContentID: contentID,
			UserID:    userID,
			Content:   content,
			Date:      time.Now().Unix() * 1000,
		})
		if err != nil {
			return "", err
		}
	} else {
		err := ReplyDB.Insert(&Reply{
			ID:       newComment,
			FatherID: fatherID,
			UserID:   userID,
			Content:  content,
			Date:     time.Now().Unix() * 1000,
		})
		if err != nil {
			return "", err
		}
	}
	return newComment, nil
}

// AddLike 点赞 1或-1
func AddLike(id string, isReply bool, num int) (err error) {
	if isReply {
		_, err = ReplyDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
	} else {
		_, err = CommentDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
	}
	return
}

// RemoveComment 删除评论
func RemoveComment(id string, isReply bool) (err error) {
	if isReply {
		_, err = ReplyDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"content": ""}})
	} else {
		_, err = CommentDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"content": ""}})
	}
	return
}

// GetCommentByContentID 获取内容指定页数的评论
func GetCommentByContentID(id string, eachNum, pageNum int) (comment []Comment) {
	err := CommentDB.Find(nil).Sort("-date").Skip(eachNum * (pageNum - 1)).Limit(eachNum).All(&comment)
	if err != nil {
		return nil
	}
	return
}

// GetReplyByCommentID 获取指定ID评论的回复
func GetReplyByCommentID(id string) (reply []Reply) {
	ReplyDB.Find(bson.M{"fatherId": id}).All(&reply)
	return
}

// GetCommentPage 获取评论数目
func GetCommentPage(id string) (count int) {
	count, err := CommentDB.Find(bson.M{"contentId": id}).Count()
	if err != nil {
		count = -1
	}
	return
}
