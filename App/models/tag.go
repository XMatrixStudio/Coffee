package models

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Tag 标签
type Tag struct {
	userID bson.ObjectId `bson:"userId"`
	Name   string        `bson:"name"`  // 标签名
	Count  int64         `bson:"count"` // 【索引】 使用该标签的资源数 用于排序
}

// TagModel 标签数据库
type TagModel struct {
	DB *mgo.Collection
}

// AddTag 增加指定用户的指定标签的数量
func (m *TagModel) AddTag(userID, tag string) (err error) {
	_, err = m.DB.Upsert(bson.M{"userId": bson.ObjectIdHex(userID), "name": tag}, bson.M{"$inc": bson.M{"count": 1}})
	return
}

// GetTags 获取指定用户的标签（按数量排序）
func (m *TagModel) GetTags(userID string) (tags []string, err error) {
	err = m.DB.Find(bson.M{"userId": bson.ObjectIdHex(userID)}).Sort("count").All(&tags)
	return
}
