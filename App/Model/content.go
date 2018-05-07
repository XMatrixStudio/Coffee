package model

import (
	"time"

	content "github.com/XMatrixStudio/Coffee/App/Model/Content"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Content 内容
type Content struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`        // 内容名字
	Detail      string        `bson:"detail"`      // 详情介绍
	OwnID       string        `bson:"ownId"`       // 作者ID
	PublishDate int64         `bson:"publishDate"` // 发布日期
	EditDate    int64         `bson:"editDate"`    // 修改日期
	LikeNum     int64         `bson:"likeNum"`     // 点赞人数
	CommentNum  int64         `bson:"commentNum"`  // 评论次数
	ReadNum     int64         `bson:"readNum"`     // 阅读次数
	Top         bool          `bson:"top"`         // 是否置顶
	Public      bool          `bson:"public"`      // 是否公开
	Comment     bool          `bson:"comment"`     // 是否开放评论
	Type        string        `bson:"type"`        // 类型， "Movie", "Data", "Picture"， "Doc", "App", "Log", "Daily"
	Local       bool          `bson:"local"`       // 是否本地资源
	File        content.File  `bson:"file"`        // 文件系统
	Image       string        `bson:"image"`       // 预览图URL
	Tag         []string      `bson:"tag"`         // 标签
}

// ContentDB 内容数据库
var ContentDB *mgo.Collection

// AddContent 增加内容
func AddContent(name, detail, contentType, userID string, isComment bool) (bson.ObjectId, error) {
	newContent := bson.NewObjectId()
	err := ContentDB.Insert(&Content{
		ID:          newContent,
		Name:        name,
		Detail:      detail,
		OwnID:       userID,
		PublishDate: time.Now().Unix() * 1000,
		EditDate:    time.Now().Unix() * 1000,
		LikeNum:     0,
		CommentNum:  0,
		ReadNum:     0,
		Top:         false,
		Public:      true,
		Comment:     isComment,
		Type:        contentType,
		Local:       false,
	})
	if err != nil {
		return "", err
	}
	return newContent, nil
}

// EditContent 更新修改日期
func EditContent(id string) error {
	err := ContentDB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"editDate": time.Now().Unix() * 1000}})
	return err
}

// AddNum 增加一个或减少一个阅读("readNum")/点赞("likeNum")/评论数("commentNum")
func AddNum(id, name string, num int) error {
	err := ContentDB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{name: num}})
	return err
}

// SetStatus 设置置顶("top")/评论("comment")/锁定("lock")状态
func SetStatus(id, name string, status bool) error {
	err := ContentDB.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{name: status}})
	return err
}

// RemoveContent 删除内容
func RemoveContent(id string) error {
	err := ContentDB.RemoveId(bson.ObjectIdHex(id))
	return err
}

// GetContentByID 根据ID查询内容
func GetContentByID(id string) *Content {
	content := new(Content)
	err := ContentDB.FindId(id).One(&content)
	if err != nil {
		return nil
	}
	return content
}

// GetContentByOwn 根据作者ID查询内容
func GetContentByOwn(id string) []Content {
	var content []Content
	err := ContentDB.Find(bson.M{"ownId": id}).All(&content)
	if err != nil {
		return nil
	}
	return content
}

// GetPageContent 获取内容指定分页内容集合
func GetPageContent(eachNum, pageNum int) []Content {
	var content []Content
	err := ContentDB.Find(nil).Sort("-editDate").Skip(eachNum * (pageNum - 1)).Limit(eachNum).All(&content)
	if err != nil {
		return nil
	}
	return content
}
