package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"errors"
)

// Gather 集合
type Gather struct {
	ID  bson.ObjectId `bson:"_id"`
	MID string        `bson:"id"`  // MID 【索引】(ID可以是内容，可以为评论，甚至可以是用户(关注系统))
	IDs []string      `bson:"ids"` // ID集合
}

// GatherModel 集合数据库
type GatherModel struct {
	FollowerDB *mgo.Collection
	FollowingDB *mgo.Collection
	// ContentLikeDB 点赞[内容]的[用户ID]集合
	// 点赞某内容的用户， 用于详细页面显示，属于冗余数据
	ContentLikeDB *mgo.Collection
	// UserLikeDB [用户]点赞[内容或评论的ID]的集合
	UserLikeDB *mgo.Collection
}

func (m *GatherModel) RemoveAllByID(id string) {
	m.ContentLikeDB.Remove(bson.M{"id": id})
	m.UserLikeDB.Upsert(nil, bson.M{"$pull": bson.M{"ids": id}})
}


// AddLikeToContent 增加Like到内容里面
func (m *GatherModel) AddLikeToContent(contentID, userID string) error {
	info, err := m.ContentLikeDB.Upsert(bson.M{"id": contentID}, bson.M{"$addToSet": bson.M{"ids": userID}})
	if info.Matched == 1 && info.Updated == 0 {
		return errors.New("exist")
	}
	return err
}

// AddLikeToUser
func (m *GatherModel) AddLikeToUser(contentID, userID string) error {
	info, err := m.UserLikeDB.Upsert(bson.M{"id": userID}, bson.M{"$addToSet": bson.M{"ids": contentID}})
	if info.Matched == 1 && info.Updated == 0 {
		return errors.New("exist")
	}
	return err
}

// RemoveLikeFromContent 取消点赞内容
func (m *GatherModel) RemoveLikeFromContent(contentID, userID string) error {
	info, err := m.ContentLikeDB.Upsert(bson.M{"id": contentID}, bson.M{"$pull": bson.M{"ids": userID}})
	if info.Updated == 0 {
		return errors.New("not found")
	}
	return err
}

// RemoveLikeFromUser
func (m *GatherModel) RemoveLikeFromUser(commentID, userID string) error {
	info, err := m.UserLikeDB.Upsert(bson.M{"id": userID}, bson.M{"$pull": bson.M{"ids": commentID}})
	if info.Updated == 0 {
		return errors.New("not found")
	}
	return err
}


// GetUserLikes 获取用户点赞的集合
func (m *GatherModel) GetUserLikes(userID string) ([]string, error) {
	var gather Gather
	err := m.UserLikeDB.Find(bson.M{"id": userID}).One(&gather)
	if err != nil {
		return nil, err
	}
	return gather.IDs, nil
}
