package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Achievement 日志记录
type Achievement struct {
	ID              bson.ObjectId `bson:"_id"`
	AchievementType string        `bson:"achievementType"` // 成就类型

}

// AchievementModel 日志数据库
type AchievementModel struct {
	DB *mgo.Collection
}
