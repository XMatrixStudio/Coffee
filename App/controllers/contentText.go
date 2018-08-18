package controllers

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"html/template"
	"github.com/XMatrixStudio/Coffee/App/services"
)

// GetTextsBy GET /content/texts/{id} 获取指定用户的所有Text类型的内容
func (c *ContentController) GetTextsBy(id string) (res ContentsRes) {
	isOwn := false
	if id == "self" {
		if c.Session.Get("id") == nil {
			res.State = "not_login"
			return
		}
		id = c.Session.GetString("id")
		isOwn = true
	} else if !bson.IsObjectIdHex(id) {
		res.State = "error_id"
		return
	}
	res.State = "success"
	res.Data = c.Service.GetTextByUser(id, !isOwn)
	return
}

// PostText POST /content/text 增加文本内容
func (c *ContentController) PostText() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := services.ContentData{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Title == "" || req.Detail == "" {
		res.State = "bad_request"
		return
	}
	req.Detail = template.HTMLEscapeString(req.Detail)
	err = c.Service.AddText(c.Session.GetString("id"), req)
	if err != nil {
		res.State = "error"
		return
	}
	res.State = "success"
	return
}

// PatchTextBy PATCH /content/all/{contentID} 修改指定文本内容
func (c *ContentController) PatchAllBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	fmt.Println(id)
	if id == "" {
		res.State = "bad_request"
		return
	}
	req := services.ContentData{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Title == "" || req.Detail == "" {
		res.State = "bad_request"
		res.Data = err.Error()
		return
	}
	req.Detail = template.HTMLEscapeString(req.Detail)
	err = c.Service.PatchContentByID(id, req)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
		return
	}
	res.State = "success"
	return
}
