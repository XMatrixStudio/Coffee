package models

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

// Model ...
type Model struct {
	Config       Mongo
	DB           *mgo.Database
	User         UserModel
	Comment      CommentModel
	Content      ContentModel
	Gather       GatherModel
	Notification NotificationModel
	Tag          TagModel
}

// InitMongo 初始化数据库
func (m *Model) InitMongo(conf Mongo) error {
	m.Config = conf
	if m.DB != nil {
		m.DB.Session.Close()
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
	m.DB = session.DB(conf.Name)
	m.User.DB = m.DB.C("users")
	m.Content.DB = m.DB.C("contents")
	m.Tag.DB = m.DB.C("tags")
	m.Notification.DB = m.DB.C("notifications")
	m.Comment.CommentDB = m.DB.C("comments")
	m.Comment.ReplyDB = m.DB.C("reply")
	m.Gather.ContentLikeDB = m.DB.C("like")
	m.Gather.UserLikeDB = m.DB.C("userLike")
	m.Gather.FollowingDB = m.DB.C("following")
	m.Gather.FollowerDB = m.DB.C("follower")
	log.Printf("MongoDB Connect Success!")
	return nil
}

// NewModel 初始化Model
func NewModel(c Mongo) (*Model, error) {
	model := new(Model)
	err := model.InitMongo(c)
	return model, err
}
