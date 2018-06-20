package models

import (
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Content 内容
/* 类型列表
- Movie 电影 - Movie
- Album 相册 - Album
- App 应用 - App
 - Android
 - Windows
 - Linux
- Game 游戏 - App
 - Android
 - Windows
- Files 文件 - Files
 - Docs
 - Music
 - PDF
 - Backup
 - Code
 - Other
- Daily 日记/随笔
*/
type Content struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string        `bson:"name"`        // 内容名字
	Detail      string        `bson:"detail"`      // 详情介绍
	OwnID       bson.ObjectId `bson:"ownId"`       // 作者ID [索引]
	PublishDate int64         `bson:"publishDate"` // 发布日期
	EditDate    int64         `bson:"editDate"`    // 修改日期
	LikeNum     int64         `bson:"likeNum"`     // 点赞人数
	CommentNum  int64         `bson:"commentNum"`  // 评论次数
	Public      bool          `bson:"public"`      // 是否公开
	Native      bool          `bson:"native"`      // 是否本地资源
	Type        string        `bson:"type"`        // 类型， "Movie", "Data", "Album"， "Docs", "App", "Daily"
	SubType     string        `bson:"subType"`     // 子类型 (如 "app"-"android", "windows")
	Remarks     string        `bson:"remark"`      // 备注
	Tag         []string      `bson:"tag"`         // 标签（ObjectId）
	Image       []Image       `bson:"image"`       // 预览图
	Files       []File        `bson:"append"`      // 文件集合 (可以用于存储电影字幕，软件附件等)
	Movie       Movie         `bson:"movie"`       // Movie类型专属
	Album       Album         `bson:"album"`       // Album类型专属
	App         App           `bson:"app"`         // App/Game类型专属
}

// ContentModel 内容数据库
type ContentModel struct {
	DB *mgo.Collection
}

// AddContent 增加内容
func (m *ContentModel) AddContent(name, detail, userID, contentType string) (bson.ObjectId, error) {
	newContent := bson.NewObjectId()
	err := m.DB.Insert(&Content{
		ID:          newContent,
		Name:        name,
		Detail:      detail,
		OwnID:       bson.ObjectIdHex(userID),
		PublishDate: time.Now().Unix() * 1000,
		EditDate:    time.Now().Unix() * 1000,
		LikeNum:     0,
		CommentNum:  0,
		Public:      true,
		Type:        contentType,
		Native:      false,
	})
	if err != nil {
		return "", err
	}
	return newContent, nil
}

// RemoveContent 删除内容
func (m *ContentModel) RemoveContent(id string) (err error) {
	err = m.DB.RemoveId(bson.ObjectIdHex(id))
	return
}

// GetContentByID 根据ID查询内容
func (m *ContentModel) GetContentByID(id string) *Content {
	content := new(Content)
	err := m.DB.FindId(id).One(&content)
	if err != nil {
		return nil
	}
	return content
}

// UpdateByID 根据ID更新内容
func (m *ContentModel) UpdateByID(id string, data Content) (err error) {
	err = m.DB.UpdateId(bson.ObjectIdHex(id), &data)
	return
}

// GetContentByOwn 根据作者ID查询内容
func (m *ContentModel) GetContentByOwn(ownID string) []Content {
	var content []Content
	err := m.DB.Find(bson.M{"ownId": bson.ObjectIdHex(ownID)}).All(&content)
	if err != nil {
		return nil
	}
	return content
}

func (m *ContentModel) GetCountByOwn(ownID string) (count int, err error) {
	count, err = m.DB.Find(bson.M{"ownId": bson.ObjectIdHex(ownID)}).Count()
	return
}

// GetPageContent 获取内容指定分页内容集合
func (m *ContentModel) GetPageContent(ownID, contentType, subType string, eachNum, pageNum int) []Content {
	var content []Content
	err := m.DB.Find(bson.M{"ownId": bson.ObjectIdHex(ownID)}).Sort("-editDate").Skip(eachNum * (pageNum - 1)).Limit(eachNum).All(&content)
	if err != nil {
		return nil
	}
	return content
}
