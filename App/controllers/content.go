package controllers

import (
	"strconv"

	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
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

// ContentsRes 内容集合回复
type ContentsRes struct {
	State string
	Data  []models.Content
}

// GetDetailBy GET /content/detail/{contentID} 获取指定内容
func (c *ContentController) GetDetailBy(id string) (res ContentRes) {
	if !bson.IsObjectIdHex(id) {
		res.State = StatusBadReq
		return
	}
	content, user, err := c.Service.GetContentAndUser(id)
	if err != nil {
		res.State = err.Error()
		return
	}
	if content.Public == false && c.Session.Get("id") == nil && c.Session.GetString("id") != content.OwnID.Hex() {
		res.State = StatusNotLogin
		return
	}
	res.Data = content
	res.User = user
	res.State = StatusSuccess
	return
}

// PublishRes 公共内容返回值
type PublishRes struct {
	State string
	Data  []services.PublishData
}

// GetPublic GET /content/public 获取公共内容
func (c *ContentController) GetPublic() (res PublishRes) {
	page, err := strconv.Atoi(c.Ctx.FormValue("page"))
	if err != nil {
		res.State = StatusBadReq
		return
	}
	eachPage, err := strconv.Atoi(c.Ctx.FormValue("eachPage"))
	if err != nil {
		res.State = StatusBadReq
		return
	}
	res.State = StatusSuccess
	res.Data = c.Service.GetPublicContents(page, eachPage)
	return
}

// DeleteBy DELETE /content/{contentID} 删除指定内容
func (c *ContentController) DeleteBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = StatusBadReq
		return
	}
	err := c.Service.DeleteContentByID(id, c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = StatusSuccess
	return
}
