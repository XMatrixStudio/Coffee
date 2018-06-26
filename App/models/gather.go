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

// RemoveAllByID 删除与指定内容有关的数据
func (m *GatherModel) RemoveAllByID(id string) {
	m.ContentLikeDB.Remove(bson.M{"id": id})
	m.UserLikeDB.Upsert(nil, bson.M{"$pull": bson.M{"ids": id}})
}

// removeFromGather 从集合中删除内容
func (m *GatherModel) removeFromGather(DB *mgo.Collection, mID, cID string) error {
	info, err := DB.Upsert(bson.M{"id": mID}, bson.M{"$pull": bson.M{"ids": cID}})
	if info.Updated == 0 {
		return errors.New("not found")
	}
	return err
}

// addTOGather 向集合中添加内容
func (m *GatherModel) addToGather(DB *mgo.Collection, mID, cID string) error {
	info, err := DB.Upsert(bson.M{"id": mID}, bson.M{"$addToSet": bson.M{"ids": cID}})
	if info.Matched == 1 && info.Updated == 0 {
		return errors.New("exist")
	}
	return err
}

func (m *GatherModel) getGather(DB *mgo.Collection, mID string) (IDs []string, err error) {
	var gather Gather
	if  DB.Find(bson.M{"id": mID}).One(&gather) != nil {
		return
	}
	IDs = gather.IDs
	return
}

// AddLikeToContent 增加Like到内容里面
func (m *GatherModel) AddLikeToContent(contentID, userID string) error {
	return m.addToGather(m.ContentLikeDB, contentID, userID)
}

// AddLikeToUser
func (m *GatherModel) AddLikeToUser(contentID, userID string) error {
	return m.addToGather(m.UserLikeDB, userID, contentID)
}

// AddUserFollowing 增加用户的关注列表
func (m *GatherModel) AddUserFollowing(userID, followingID string) error {
	return m.addToGather(m.FollowingDB, userID, followingID)
}

// AddUserFollower 增加用户的被关注列表
func (m *GatherModel) AddUserFollower(userID, followerID string) error {
	return m.addToGather(m.FollowerDB, userID, followerID)
}

// RemoveLikeFromContent 取消点赞内容
func (m *GatherModel) RemoveLikeFromContent(contentID, userID string) error {
	return m.removeFromGather(m.ContentLikeDB, contentID, userID)
}

// RemoveLikeFromUser
func (m *GatherModel) RemoveLikeFromUser(commentID, userID string) error {
	return m.removeFromGather(m.UserLikeDB, userID, commentID)
}

// RemoveFollowing
func (m *GatherModel) RemoveFollowing(userID, FollowingID string) error {
	return m.removeFromGather(m.FollowingDB, userID, FollowingID)
}

// RemoveFollower
func (m *GatherModel) RemoveFollower(userID, FollowerID string) error {
	return m.removeFromGather(m.FollowerDB, userID, FollowerID)
}

// GetUserLikes 获取用户点赞的集合
func (m *GatherModel) GetUserLikes(userID string) (IDs []string,err error) {
	IDs, err = m.getGather(m.UserLikeDB, userID)
	return
}

// GetContentLikes 获取文档点赞的用户
func (m *GatherModel) GetContentLikes(contentID string) (IDs []string,err error)  {
	IDs, err = m.getGather(m.UserLikeDB, contentID)
	return
}

// GetFollowing
func (m *GatherModel) GetFollowing(userID string) (IDs []string,err error) {
	IDs, err = m.getGather(m.FollowingDB, userID)
	return
}

// GetFollower
func (m *GatherModel) GetFollower(userID string) (IDs []string,err error) {
	IDs, err = m.getGather(m.FollowerDB, userID)
	return
}