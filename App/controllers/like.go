package controllers

import (
	"github.com/kataras/iris"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris/sessions"
)

// LikeController Like
type LikeController struct {
	Ctx     iris.Context
	Service services.LikeService
	Session *sessions.Session
}