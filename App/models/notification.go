package models

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Notification 用户通知
type Notification struct {
	ID            bson.ObjectId        `bson:"_id"`
	UserID        bson.ObjectId        `bson:"userId"`        // 用户ID 【索引】
	Notifications []NotificationDetail `bson:"notifications"` // 通知集合
}

// 通知类型
const (
	TypeSystem = "system"
	TypeLike   = "like"
	TypeReply  = "reply"
)

// NotificationDetail 通知详情
type NotificationDetail struct {
	ID         bson.ObjectId `bson:"_id"`
	CreateTime int64         `bson:"time"`
	Content    string        `bson:"content"`  // 通知内容
	SourceID   string        `bson:"sourceId"` // 源ID （点赞人）
	TargetID   string        `bson:"targetId"` // 目标ID （点赞文章）
	Read       bool          `bson:"read"`     // 是否已读
	Type       string        `bson:"type"`     // 类型： "system", "like", "reply"...
}

// NotificationModel 通知数据库
type NotificationModel struct {
	DB *mgo.Collection
}

// AddNotification 添加一条通知 类型:"system", "like", "reply" ...
func (m *NotificationModel) AddNotification(content, user, sourceID, targetID, notificationType string) error {
	newNotification := &NotificationDetail{
		ID:         bson.NewObjectId(),
		CreateTime: time.Now().Unix() * 1000,
		Content:    content,
		Read:       false,
		SourceID:   sourceID,
		TargetID:   targetID,
		Type:       notificationType,
	}
	_, err := m.DB.Upsert(
		bson.M{"userId": bson.ObjectIdHex(user)},
		bson.M{"$push": bson.M{"notifications": &newNotification}})
	return err
}

// ReadANotification 标记通知
func (m *NotificationModel) ReadANotification(user, id string, status bool) error {
	return m.DB.Update(
		bson.M{"userId": bson.ObjectIdHex(user), "notifications._id": bson.ObjectIdHex(id)},
		bson.M{"$set": bson.M{"notifications.$.read": status}})
}

// RemoveANotification 删除通知
func (m *NotificationModel) RemoveANotification(user, id string) error {
	return m.DB.Update(
		bson.M{"userId": bson.ObjectIdHex(user)},
		bson.M{"$pull": bson.M{"notifications": bson.M{"_id": bson.ObjectIdHex(id)}}})
}

// GetNotificationsByUser 获取用户所有通知
func (m *NotificationModel) GetNotificationsByUser(user string) ([]NotificationDetail, error) {
	var notification Notification
	err := m.DB.Find(bson.M{"userId": bson.ObjectIdHex(user)}).One(&notification)
	return notification.Notifications, err
}

// GetUnreadCountByUser 获取用户未读通知数量
func (m *NotificationModel) GetUnreadCountByUser(userID string) (count int, err error) {
	count, err = m.DB.Find(bson.M{"userId": bson.ObjectIdHex(userID), "notifications.read": false}).Count()
	return
}

func (m *NotificationModel) RemoveUnread(user, sid, tid string) error {
	return m.DB.Update(
		bson.M{"userId": bson.ObjectIdHex(user)},
		bson.M{"$pull": bson.M{"notifications": bson.M{"notifications.targetId": tid, "notifications.sourceId": sid, "notifications.read": false}}})
}
