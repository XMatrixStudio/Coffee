package models

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Notification 用户通知
type Notification struct {
	ID            bson.ObjectId        `bson:"_id"`
	UserID        string               `bson:"userId"`        // 用户ID 【索引】
	Notifications []NotificationDetail `bson:"notifications"` // 通知集合
}

// NotificationDetail 通知详情
type NotificationDetail struct {
	ID      bson.ObjectId `bson:"_id"`
	Content string        `bson:"content"` // 通知内容
	Read    bool          `bson:"read"`    // 是否以读
	Type    string        `bson:"type"`    // 类型： "system", "like", "reply"...
}

// NotificationModel 通知数据库
type NotificationModel struct {
	DB *mgo.Collection
}

// InitNotification 初始化用户通知
func (m *NotificationModel) InitNotification(user string) error {
	newNotification := bson.NewObjectId()
	err := m.DB.Insert(&Notification{
		ID:     newNotification,
		UserID: user,
	})
	return err
}

// AddNotification 添加一条通知 类型:"system", "like", "reply" ...
func (m *NotificationModel) AddNotification(content, user, notificationType string) error {
	newNotification := &NotificationDetail{
		ID:      bson.NewObjectId(),
		Content: content,
		Read:    false,
		Type:    notificationType,
	}
	err := m.DB.Update(
		bson.M{"userId": bson.ObjectIdHex(user)},
		bson.M{"$push": bson.M{"notifications": &newNotification}})
	return err
}

// ReadANotification 标记通知
func (m *NotificationModel) ReadANotification(user, id string, status bool) error {
	err := m.DB.Update(
		bson.M{"userId": bson.ObjectIdHex(user), "notifications._id": id},
		bson.M{"$set": bson.M{"notifications.$.read": status}})
	return err
}

// RemoveANotification 删除通知
func (m *NotificationModel) RemoveANotification(user, id string) error {
	err := m.DB.Update(
		bson.M{"userId": bson.ObjectIdHex(user), "notifications._id": id},
		bson.M{"$pull": bson.M{"notifications._id": id}})
	return err
}

// GetNotificationsByUser 获取用户所有通知
func (m *NotificationModel) GetNotificationsByUser(user string) ([]NotificationDetail, error) {
	var notifications []NotificationDetail
	err := m.DB.Find(nil).All(notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
