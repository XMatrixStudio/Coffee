package controllers

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

// GetTexts GET /content/texts 获取指定用户的所有内容
func (c *ContentController) GetTexts() (res ContentsRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	res.State = "success"
	res.Data = c.Service.GetTextByUser(c.Session.GetString("id"), false)
	return
}

func (c *ContentController) GetTextsBy(id string) (res ContentsRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = "error_req"
		return
	}
	res.State = "success"
	res.Data = c.Service.GetTextByUser(c.Session.GetString("id"), true)
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
	err = c.Service.AddText(c.Session.GetString("id"), req.Title, req.Content, req.IsPublic, req.Tags)
	if err != nil {
		res.State = "error"
		return
	}
	res.State = "success"
	return
}

// PatchTextBy PATCH /content/text/{contentID} 修改指定文本内容
func (c *ContentController) PatchTextBy(id string) (res CommonRes) {
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
	err = c.Service.PatchContentByID(id, req.Title, req.Content, req.Tags, req.IsPublic)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
		return
	}
	res.State = "success"
	return
}
