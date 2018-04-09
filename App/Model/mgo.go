package model

import (
	"github.com/globalsign/mgo"
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
	return nil
}
