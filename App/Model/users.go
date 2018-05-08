package model

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Users 用户基本信息
/*
Class 用户类型
- -2： 黑名单用户
 - 无任何权限(无法登陆)
- -1： 受限用户
 - 浏览
- 0： 普通用户
 - 浏览
 - 下载
 - 评论
 - 点赞
 - 发布内容
- 1： 认证用户
 - 继承普通用户权限
 - 上传图片
 - 上传文件
- 10： VIP用户
 - 继承认证用户权限
 - VIP功能
- 50： 管理员
 - 继承VIP用户权限
 - 调整VIP及以下的用户级别
 - 调整VIP及以下的用户存储库大小限制
 - 删除任意公开评论
 - 删除任意公开内容(调整为私有)
- 99： 超级管理员
 - 继承管理员权限
 - 调整任意用户级别
 - 修改系统设置

Size 存储库分配
- 认证用户
 - 最大32G， 单个8G
- VIP用户 - 高级用户
 - 最大512G， 单个128G
- VIP用户 - 超级用户
 - 最大2048G， 单个256G
- 管理员
 - 无上限
- 超级管理员
 - 无上限
*/
type Users struct {
	ID         bson.ObjectId `bson:"_id"`        // 用户ID
	Name       string        `bson:"name"`       // 用户唯一名字
	Class      int           `bson:"class"`      // 用户类型
	Info       UserInfo      `bson:"info"`       // 用户个性信息
	LikeNum    int64         `bson:"likeNum"`    // 收到的点赞数
	Token      string        `bson:"token"`      // Violet 访问令牌
	MaxSize    int64         `bson:"maxSize"`    // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize   int64         `bson:"usedSize"`   // 存储库已用大小 单位为KB
	SingleSize int64         `bson:"singleSize"` // 单个资源最大上限 -1为无上限
	FilesClass []string      `bson:"filesClass"` // 文件分类
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
		ID:         newUser,
		Name:       "user_" + string(newUser),
		Class:      0,
		FilesClass: []string{"文档", "图书", "音乐", "代码", "备份", "其他"}, // 默认分类
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
