package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/XMatrixStudio/Coffee/App/models"
)

// ContentController 内容
type ContentController struct {
	Ctx     iris.Context
	Service services.ContentService
	Session *sessions.Session
}

// ContentRes 内容回复
type ContentRes struct {
	State string
	Data  models.Content
	User  services.UserBaseInfo
}

func (c *ContentController) GetContentBy(id string) (res ContentRes) {
	content, err := c.Service.GetContentByID(id)
	if err != nil {
		res.State = err.Error()
		return
	}
	if content.Public == false && c.Session.Get("id") == nil && c.Session.GetString("id") != content.OwnID.Hex() {
		res.State = "not_login"
		return
	}
	res.Data = content
	res.User = c.Service.GetUserBaseInfo(content.OwnID.Hex())
	res.State = "success"
	return
}
