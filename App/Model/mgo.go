package model

import (
	"log"

	"github.com/globalsign/mgo"
)

// Mongo 数据库配置
type Mongo struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Name     string `yaml:"Name"`
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
	ReplyDB = DB.C("reply")
	ContentLikeDB = DB.C("like")
	CommentLikeDB = DB.C("commentLike")
	UserLikeDB = DB.C("userLike")
	NotificationDB = DB.C("notifications")
	TagDB = DB.C("tags")
	log.Printf("MongoDB Connect Success!")
	return nil
}
