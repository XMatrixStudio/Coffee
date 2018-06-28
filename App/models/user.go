package models

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

// 用户类型
const (
	ClassBlackUser  int = iota
	ClassLimitUser
	ClassNormalUser
	ClassVerifyUser
	ClassVIPUser
	ClassSVIPUser
	ClassAdmin
	ClassSAdmin
)

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
 - 最大8G， 单个2G
- VIP用户 - 高级用户
 - 最大128G， 单个8G
- VIP用户 - 超级用户
 - 最大1024G， 单个128G
- 管理员
 - 无上限
- 超级管理员
 - 无上限
*/
// User 用户基本信息
type User struct {
	ID       bson.ObjectId `bson:"_id"`   // 用户ID
	VioletID bson.ObjectId `bson:"vid"`   // VioletID
	Email    string        `bson:"email"` // 用户唯一邮箱
	Class    int           `bson:"class"` // 用户类型
	Info     UserInfo      `bson:"info"`  // 用户个性信息
	Token    string        `bson:"token"` // Violet 访问令牌

	MaxSize    int64 `bson:"maxSize"`    // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize   int64 `bson:"usedSize"`   // 存储库已用大小 单位为KB
	SingleSize int64 `bson:"singleSize"` // 单个资源最大上限 -1为无上限

	FilesClass []string `bson:"filesClass"` // 文件分类


	LikeCount      int64     `bson:"likeCount"`      // 被点赞数
	MaxLikeCount   int64     `bson:"maxLikeCount"`   // 最大被点赞数
	CommentTime    time.Time `bson:"commentTime"`    // 最后一次评论时间
	ContentTime    time.Time `bson:"contentTime"`    // 最后一次发布内容时间
	FollowerCount  int64     `bson"followerCount"`   // 被关注数目
	FollowingCount int64     `bson:"followingCount"` // 关注数目
	Exp            int64     `bson:"exp"`            // 经验
}

// 性别
const (
	GenderMan     int = iota
	GenderWoman
	GenderUnknown
)

// UserInfo 用户个性信息
type UserInfo struct {
	Name   string `bson:"name"`   // 用户昵称
	Avatar string `bson:"avatar"` // 头像URL
	Bio    string `bson:"bio"`    // 个人简介
	Gender int    `bson:"gender"` // 性别
}

// UserModel 用户数据库
type UserModel struct {
	DB *mgo.Collection
}

// AddUser 添加用户
func (m *UserModel) AddUser(vID, token, email, name, avatar, bio string, gender int) (newUser bson.ObjectId, err error) {
	if !bson.IsObjectIdHex(vID) {
		err = errors.New("not_id")
		return
	}
	newUser = bson.NewObjectId()
	err = m.DB.Insert(&User{
		ID:         newUser,
		VioletID:   bson.ObjectIdHex(vID),
		Email:      email,
		Class:      1,
		FilesClass: []string{"文档", "图书", "音乐", "代码", "备份", "其他"}, // 默认分类
		Info: UserInfo{
			Name:   name,
			Avatar: avatar,
			Bio:    bio,
			Gender: gender,
		},
		Token:      token,
		MaxSize:    8388608,
		SingleSize: 2097152,
	})
	return
}

// SetUserInfo 设置用户信息
func (m *UserModel) SetUserInfo(id string, info UserInfo) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"info": info}})
}

// SetUserName 设置用户名
func (m *UserModel) SetUserName(id, name string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"info.name": name}})
}

// GetUserByID 根据ID查询用户
func (m *UserModel) GetUserByID(id string) (user User, err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New("not_id")
		return
	}
	err = m.DB.FindId(bson.ObjectIdHex(id)).One(&user)
	return
}

// GetUserByVID 根据VioletID查询用户
func (m *UserModel) GetUserByVID(id string) (user User, err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New("not_id")
		return
	}
	err = m.DB.Find(bson.M{"vid": bson.ObjectIdHex(id)}).One(&user)
	return
}

// SetUserToken 设置Token
func (m *UserModel) SetUserToken(id, token string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"token": token}})
}

// AddLikeNum 增加或减少点赞数
func (m *UserModel) AddLikeNum(id string, add bool) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	num := -1
	if add {
		num = 1
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
}

// AddContentCount 增加内容数
func (m *UserModel) AddContentCount(id string, add bool) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	num := -1
	if add {
		num = 1
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"contentCount": num}})
}

// AddUsedSize 增加已用大小
func (m *UserModel) SetUsedSize(id string, size int64) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"usedSize": size}})
}

// SetUserClass 设置用户类型
func (m *UserModel) SetUserClass(id, class string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"class": class}})
}

// SetSize 设置用户容量大小
func (m *UserModel) SetSize(id string, maxSize, singleSize int64) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"maxSize": maxSize, "singleSize": singleSize}})
}

// AddFilesClass 增加文件分类
func (m *UserModel) AddFilesClass(id string, filesClass string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$push": bson.M{"filesClass": filesClass}})
}

// DeleteFilesClass 减少文件分类
func (m *UserModel) DeleteFilesClass(id string, filesClass string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$pull": bson.M{"filesClass": filesClass}})
}
