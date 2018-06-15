package controllers

import (
	"github.com/XMatrixStudio/Coffee/app/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

type UsersController struct {
	Ctx     iris.Context
	Service services.UserService
	Session *sessions.Session
}

func (c *UsersController) GetTest() string {
	return "OK"
}
