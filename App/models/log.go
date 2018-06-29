package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// LogRecord 日志记录
type LogRecord struct {
	ID      bson.ObjectId `bson:"_id"`
	LogType string        `bson:"logType"` // 日志类型

}

// LogModel 日志数据库
type LogModel struct {
	DB *mgo.Collection
}
