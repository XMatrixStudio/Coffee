package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
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
	Data  []models.Content
}

// GetMyContent GET /MyContent 获取指定用户的所有内容
func (c *ContentController) GetMyContent() (res ContentRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	res.Data = c.Service.GetContentByOwn(c.Session.GetString("id"))
	return
}

type textDataReq struct {
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	IsPublic bool     `json:"isPublic"`
	Tags     []string `json:"tags`
}

// PostText POST /text 增加文本内容
func (c *ContentController) PostText() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := textDataReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = "bad_request"
		return
	}
	err = c.Service.AddText(c.Session.GetString("id"), req.Title, req.Text, req.IsPublic, req.Tags)
	if err != nil {
		res.State = "error"
		return
	}
	res.State = "success"
	return
}
