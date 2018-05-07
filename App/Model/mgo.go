package model

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

// Mongo 数据库配置
type Mongo struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// DB 数据库连接
var DB *mgo.Database

// InitMongo 初始化数据库
func InitMongo(conf Mongo) error {
	if DB != nil {
		DB.Session.Close()
	}
	session, err := mgo.Dial(
		"mongodb://" +
			conf.User +
			":" + conf.Password +
			"@" + conf.Host +
			":" + conf.Port +
			"/" + conf.Name)
	if err != nil {
		return err
	}
	DB = session.DB(conf.Name)
	UserDB = DB.C("users")
	ContentDB = DB.C("contents")
	CommentDB = DB.C("comments")
	ContentLikeDB = DB.C("like")
	CommentLikeDB = DB.C("commentLike")
	UserLikeDB = DB.C("userLike")
	NotificationDB = DB.C("notifications")
	TagDB = DB.C("tags")
	fmt.Println("MongoDB Connect Success!")
	return nil
}
