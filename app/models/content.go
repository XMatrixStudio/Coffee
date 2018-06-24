package models

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"errors"
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

const (
	TypeText = "Text"
	TypePic  = "Pic"
	TypeFilm = "Film"
	TypeApp  = "App"
	TypeGame = "Game"
	TypeDoc  = "Doc"
)

// ContentModel 内容数据库
type ContentModel struct {
	DB *mgo.Collection
}

// AddContent 增加内容
func (m *ContentModel) AddContent(content Content) (bson.ObjectId, error) {
	content.ID = bson.NewObjectId()
	content.PublishDate = time.Now().Unix() * 1000
	content.EditDate = time.Now().Unix() * 1000
	content.LikeNum = 0
	content.CommentNum = 0
	err := m.DB.Insert(content)
	return content.ID, err
}

// RemoveContent 删除内容
func (m *ContentModel) RemoveContent(id string) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	err = m.DB.RemoveId(bson.ObjectIdHex(id))
	return
}

// GetContentByID 根据ID查询内容
func (m *ContentModel) GetContentByID(id string) (content Content,err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New("not_id")
		return
	}
	err = m.DB.FindId(bson.ObjectIdHex(id)).One(&content)
	return
}

// UpdateByID 根据ID更新内容
func (m *ContentModel) UpdateByID(id string, data Content) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.UpdateId(bson.ObjectIdHex(id), &data)

}

func (m *ContentModel) DeleteByID(id string) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("not_id")
	}
	return m.DB.RemoveId(bson.ObjectIdHex(id))
}

// GetContentByOwn 根据作者ID查询内容
func (m *ContentModel) GetContentByOwn(ownID string) []Content {
	if !bson.IsObjectIdHex(ownID) {
		return nil
	}
	var content []Content
	err := m.DB.Find(bson.M{"ownId": bson.ObjectIdHex(ownID)}).Sort("-editDate").All(&content)
	if err != nil {
		return nil
	}
	return content
}

// GetCountByOwn 获取公开内容数量
func (m *ContentModel) GetCountByOwn(ownID string) (count int, err error) {
	if !bson.IsObjectIdHex(ownID) {
		return 0, errors.New("not_id")
	}
	count, err = m.DB.Find(bson.M{"public": true}).Count()
	return
}

// GetPageContent 获取内容指定分页内容集合
func (m *ContentModel) GetPageContent(eachNum, pageNum int) []Content {
	var content []Content
	err := m.DB.Find(bson.M{"public": true}).Sort("-editDate").Skip(eachNum * (pageNum - 1)).Limit(eachNum).All(&content)
	if err != nil {
		return nil
	}
	return content
}

func (m *ContentModel) AddLikeCount(id string, num int) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"likeNum": num}})
}

func (m *ContentModel) AddCommentCount(id string, num int) error {
	return m.DB.UpdateId(bson.ObjectIdHex(id), bson.M{"$inc": bson.M{"commentNum": num}})
}