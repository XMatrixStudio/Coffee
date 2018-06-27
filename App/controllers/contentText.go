package controllers

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"html/template"
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

// textDataReq 文本内容数据请求
type textDataReq struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	IsPublic bool     `json:"isPublic"`
	Tags     []string `json:"tags"`
}

// PostText POST /content/text 增加文本内容
func (c *ContentController) PostText() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := textDataReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Title == "" || req.Content == "" {
		res.State = "bad_request"
		return
	}
	req.Content = template.HTMLEscapeString(req.Content)
	err = c.Service.AddText(c.Session.GetString("id"), req.Title, req.Content, req.IsPublic, req.Tags)
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
	req := textDataReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Title == "" || req.Content == "" {
		res.State = "bad_request"
		res.Data = err.Error()
		return
	}
	req.Content = template.HTMLEscapeString(req.Content)
	err = c.Service.PatchContentByID(id, req.Title, req.Content, req.Tags, req.IsPublic)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
		return
	}
	res.State = "success"
	return
}
