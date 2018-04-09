package model

import (
	"github.com/XMatrixStudio/icecream/config"
	"github.com/globalsign/mgo"
)

type Mongo struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type model struct {
	DB *mgo.Database
}

func NewModel()

func InitMongo(conf Mongo) error {
	session, err := mgo.Dial(
		"mongodb://" +
			conf.User +
			":" + conf.Password +
			"@" + conf.Host +
			":" + conf.Port +
			"/" + conf.Dbname)
	if err != nil {
		return err
	}
	DB = session.DB(config.Mongo.Name)
	return nil
}
