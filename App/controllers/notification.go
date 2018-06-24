package controllers

import (
	"github.com/kataras/iris"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris/sessions"
)

// NotificationController Like
type NotificationController struct {
	Ctx     iris.Context
	Service services.LikeService
	Session *sessions.Session
}

