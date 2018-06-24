package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/XMatrixStudio/Coffee/App/models"
	"strconv"
	"github.com/globalsign/mgo/bson"
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

// ContentRes 内容回复
type ContentsRes struct {
	State string
	Data  []models.Content
}

func (c *ContentController) GetContentBy(id string) (res ContentRes) {
	if !bson.IsObjectIdHex(id) {
		res.State = "err_id"
		return
	}
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

type PublishRes struct {
	State string
	Data []services.PublishData
}

func (c *ContentController) GetPublic() (res PublishRes) {
	page, err := strconv.Atoi(c.Ctx.FormValue("page"))
	if err != nil {
		res.State = "error_page"
		return
	}
	eachPage, err := strconv.Atoi(c.Ctx.FormValue("eachPage"))
	if err != nil {
		res.State = "error_pageEach"
		return
	}
	res.State = "success"
	res.Data = c.Service.GetPublic(page, eachPage)
	return
}