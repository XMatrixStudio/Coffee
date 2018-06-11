package models

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Tag 标签
type Tag struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string        `bson:"name"`  // 标签名
	Count int64         `bson:"count"` // 【索引】 使用该标签的资源数 用于排序
}

var (
	// TagDB 标签数据库
	TagDB *mgo.Collection
)
