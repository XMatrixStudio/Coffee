package model

import (
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Comment 评论
type Comment struct {
	ID        bson.ObjectId `bson:"_id"`
	ArticleID string        `bson:"articleId"` // 文章ID 【索引】
	UserID    string        `bson:"userId"`    // 评论用户ID 【索引】
	Date      int64         `bson:"date"`      // 发布时间
	Content   string        `bson:"content"`   // 评论内容
	FatherID  string        `bson:"fatherId"`  // 父评论ID
	LikeNum   int64         `bson:"likeNum"`   // 点赞数
	Top       bool          `bson:"top"`       // 是否置顶
}

// CommentDB 评论数据库
var CommentDB *mgo.Collection

// AddComment 增加评论
func AddComment(article, user, content, fatherID string) (bson.ObjectId, error) {
	newComment := bson.NewObjectId()
	err := CommentDB.Insert(&Comment{
		ID:        newComment,
		ArticleID: article,
		UserID:    user,
		Content:   content,
		FatherID:  fatherID,
		Date:      time.Now().Unix() * 1000,
	})
	if err != nil {
		return "", err
	}
	return newComment, nil
}

// AddLike 点赞 1或-1
func AddLike(id string, num int) error {
	_, err := CommentDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
	return err
}

// SetTop 设置是否置顶
func SetTop(id string, status bool) error {
	_, err := CommentDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"top": status}})
	return err
}

// RemoveComment 删除评论
func RemoveComment(id string) error {
	_, err := CommentDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"content": ""}})
	return err
}

// GetCommentByContentID 获取内容指定页数的评论
func GetCommentByContentID(id string, eachNum, pageNum int) []Comment {
	var comment []Comment
	err := CommentDB.Find(nil).Sort("-date").Skip(eachNum * (pageNum - 1)).Limit(eachNum).All(&comment)
	if err != nil {
		return nil
	}
	return comment
}

// GetCommentPage 获取评论数目
func GetCommentPage(id string) (count int) {
	count, err := CommentDB.FindId(bson.ObjectIdHex(id)).Count()
	if err != nil {
		count = -1
	}
	return
}
