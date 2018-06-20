package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// 用户类型
const (
	ClassBlackUser int = iota
	ClassLimitUser
	ClassNormalUser
	ClassVerifyUser
	ClassVIPUser
	ClassSVIPUser
	ClassAdmin
	ClassSAdmin
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
type Users struct {
	ID           bson.ObjectId `bson:"_id"`          // 用户ID
	VioletID     bson.ObjectId `bson:"vid"`          // VioletID
	Email        string        `bson:"email"`        // 用户唯一邮箱
	Class        int           `bson:"class"`        // 用户类型
	Info         UserInfo      `bson:"info"`         // 用户个性信息
	LikeNum      int64         `bson:"likeNum"`      // 收到的点赞数
	Token        string        `bson:"token"`        // Violet 访问令牌
	MaxSize      int64         `bson:"maxSize"`      // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize     int64         `bson:"usedSize"`     // 存储库已用大小 单位为KB
	SingleSize   int64         `bson:"singleSize"`   // 单个资源最大上限 -1为无上限
	FilesClass   []string      `bson:"filesClass"`   // 文件分类
	ContentCount int64         `bson:"contentCount"` // 内容数量
}

// 性别
const (
	GenderMan int = iota
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
	newUser = bson.NewObjectId()
	err = m.DB.Insert(&Users{
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
		Token: token,
		MaxSize: 8388608,
		SingleSize: 2097152,
	})
	return
}

// SetUserInfo 设置用户信息
func (m *UserModel) SetUserInfo(id string, info UserInfo) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": info})
}

// SetUserName 设置用户名
func (m *UserModel) SetUserName(id, name string) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"name": name}})
}

// GetUserByID 根据ID查询用户
func (m *UserModel) GetUserByID(id string) (user Users, err error) {
	err = m.DB.FindId(bson.ObjectIdHex(id)).One(&user)
	return
}

// GetUserByVID 根据VioletID查询用户
func (m *UserModel) GetUserByVID(id string) (user Users, err error) {
	err = m.DB.Find(bson.M{"vid": bson.ObjectIdHex(id)}).One(&user)
	return
}

// SetUserToken 设置Token
func (m *UserModel) SetUserToken(id, token string) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"token": token}})
}

// AddLikeNum 增加或减少点赞数
func (m *UserModel) AddLikeNum(id string, add bool) error {
	num := -1
	if add {
		num = 1
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
}

// AddContentCount 增加内容数
func (m *UserModel) AddContentCount(id string, add bool) error {
	num := -1
	if add {
		num = 1
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"contentCount": num}})
}

// AddUsedSize 增加已用大小
func (m *UserModel) AddUsedSize(id string, addSize int64) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"usedSize": addSize}})
}

// SetUserClass 设置用户类型
func (m *UserModel) SetUserClass(id, class string) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"class": class}})
}

// SetSize 设置用户容量大小
func (m *UserModel) SetSize(id string, maxSize, singleSize int64) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"maxSize": maxSize, "singleSize": singleSize}})
}

// AddFilesClass 增加文件分类
func (m *UserModel) AddFilesClass(id string, filesClass string) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$push": bson.M{"filesClass": filesClass}})
}

// DeleteFilesClass 减少文件分类
func (m *UserModel) DeleteFilesClass(id string, filesClass string) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$pull": bson.M{"filesClass": filesClass}})

}
