package model

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
	OwnID       string        `bson:"ownId"`       // 作者ID [索引]
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

// ContentDB 内容数据库
var ContentDB *mgo.Collection

// AddContent 增加内容
func AddContent(name, detail, userID, contentType string) (bson.ObjectId, error) {
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
func RemoveContent(id string) (err error) {
	err = ContentDB.RemoveId(bson.ObjectIdHex(id))
	return
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

// UpdateByID 根据ID更新内容
func UpdateByID(id string, data Content) (err error) {
	err = ContentDB.UpdateId(bson.ObjectIdHex(id), &data)
	return
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
func GetPageContent(ownID, contentType, subType string, eachNum, pageNum int) []Content {
	var content []Content
	err := ContentDB.Find(nil).Sort("-editDate").Skip(eachNum * (pageNum - 1)).Limit(eachNum).All(&content)
	if err != nil {
		return nil
	}
	return content
}
