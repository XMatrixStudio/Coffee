package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Users 用户基本信息
type Users struct {
	ID         bson.ObjectId `bson:"_id"`        // 用户ID
	Name       string        `bson:"name"`       // 用户唯一名字
	Class      string        `bson:"class"`      // 用户类型 "Administrator", "Subscriber", "writer"
	Info       UserInfo      `bson:"info"`       // 用户个性信息
	LikeNum    int64         `bson:"likeNum"`    // 收到的点赞数
	Token      string        `bson:"token"`      // Violet 访问令牌
	MaxSize    int64         `bson:"maxSize"`    // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize   int64         `bson:"usedSize"`   // 存储库已用大小 单位为KB
	SingleSize int64         `bson:"singleSize"` // 单个资源最大上限
}

// UserInfo 用户个性信息
type UserInfo struct {
	Avatar   string `bson:"avatar"`   // 头像URL
	Bio      string `bson:"bio"`      // 个人简介
	Email    string `bson:"email"`    // 邮箱
	Gender   int    `bson:"gender"`   // 性别
	NikeName string `bson:"nikeName"` // 昵称
}

// UserDB 数据库连接
var UserDB *mgo.Collection

// AddUser 添加用户
func AddUser() (bson.ObjectId, error) {
	newUser := bson.NewObjectId()
	err := UserDB.Insert(&Users{
		ID:    newUser,
		Name:  "user_" + string(newUser),
		Class: "reader",
	})
	if err != nil {
		return "", err
	}
	err = InitNotification(string(newUser))
	if err != nil {
		return "", err
	}
	return newUser, nil
}

// SetUserInfo 设置用户信息
func SetUserInfo(id string, info UserInfo) error {
	data := bson.M{"$set": info}
	_, err := UserDB.UpsertId(bson.ObjectIdHex(id), data)
	return err
}

// SetUserName 设置用户名
func SetUserName(id, name string) error {
	_, err := UserDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"name": name}})
	return err
}

// SetUserClass 设置用户类型
func SetUserClass(id, class string) error {
	_, err := UserDB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"class": class}})
	return err
}

// GetUserByID 根据ID查询用户
func GetUserByID(id string) (*Users, error) {
	user := new(Users)
	err := UserDB.FindId(id).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
