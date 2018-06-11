package models

import (
	"errors"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Gather 集合
type Gather struct {
	ID  bson.ObjectId `bson:"_id"`
	MID string        `bson:"id"`  // MID 【索引】(ID可以是内容，可以为评论，甚至可以是用户(关注系统))
	IDs []string      `bson:"ids"` // ID集合
}

var (
	// ContentLikeDB 点赞[内容]的[用户ID]集合
	ContentLikeDB *mgo.Collection
	// CommentLikeDB 点赞[评论]的[用户ID]集合
	CommentLikeDB *mgo.Collection
	// UserLikeDB [用户]点赞[内容或评论的ID]的集合
	// 当用户打开评论列表或者内容列表，可以一次性获取从而避免多次查询，提高性能，为冗余数据。
	UserLikeDB *mgo.Collection
)

// AddLikeToContent 增加Like到内容里面
func AddLikeToContent(contentID, userID string) error {
	info, err := ContentLikeDB.Upsert(bson.M{"id": contentID}, bson.M{"$addToSet": bson.M{"ids": userID}})
	if info.Matched == 0 {
		return errors.New("Had exits")
	}
	if err != nil {
		return err
	}
	info, err = UserLikeDB.Upsert(bson.M{"id": userID}, bson.M{"$addToSet": bson.M{"ids": contentID}})
	if info.Matched == 0 {
		return errors.New("Had exits")
	}
	return err
}

// RemoveLikeFromContent 取消点赞内容
func RemoveLikeFromContent(contentID, userID string) error {
	err := ContentLikeDB.Update(bson.M{"id": contentID}, bson.M{"$pull": bson.M{"ids": userID}})
	if err != nil {
		return err
	}
	err = UserLikeDB.Update(bson.M{"id": userID}, bson.M{"$pull": bson.M{"ids": contentID}})
	return err
}

// AddLikeToComment 增加Like到评论里面
func AddLikeToComment(commentID, userID string) error {
	info, err := CommentLikeDB.Upsert(bson.M{"id": commentID}, bson.M{"$addToSet": bson.M{"ids": userID}})
	if info.Matched == 0 {
		return errors.New("Had exits")
	}
	if err != nil {
		return err
	}
	info, err = UserLikeDB.Upsert(bson.M{"id": userID}, bson.M{"$addToSet": bson.M{"ids": commentID}})
	if info.Matched == 0 {
		return errors.New("Had exits")
	}
	return err
}

// RemoveLikeFromComment 取消点赞评论
func RemoveLikeFromComment(commentID, userID string) error {
	err := CommentLikeDB.Update(bson.M{"id": commentID}, bson.M{"$pull": bson.M{"ids": userID}})
	if err != nil {
		return err
	}
	err = UserLikeDB.Update(bson.M{"id": userID}, bson.M{"$pull": bson.M{"ids": commentID}})
	return err
}

// GetUserLikes 获取用户点赞的集合
func GetUserLikes(userID string) ([]string, error) {
	var gather Gather
	err := UserLikeDB.Find(bson.M{"id": userID}).One(&gather)
	if err != nil {
		return nil, err
	}
	return gather.IDs, nil
}
