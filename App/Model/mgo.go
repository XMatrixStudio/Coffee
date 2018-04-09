package model

import (
	"github.com/globalsign/mgo"
)

type Mongo struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

var DB *mgo.Database

func InitMongo() error {
	session, err := mgo.Dial("mongodb://")

	session.Close()
	return err
}
