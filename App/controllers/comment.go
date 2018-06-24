package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/globalsign/mgo/bson"
)

// ContentController 内容
type CommentController struct {
	Ctx     iris.Context
	Service services.CommentService
	Session *sessions.Session
}

type GetCommentRes struct {
	State string
	Data  []services.CommentForContent
}

func (c *CommentController) GetBy(id string) (res GetCommentRes) {
	if !bson.IsObjectIdHex(id) {
		res.State = "error_id"
		return
	}
	res.State = "success"
	res.Data = c.Service.GetComment(id)
	return
}

type postCommentReq struct {
	ContentID string `json:"contentId"`
	FatherID  string `json:"fatherId"`
	Content   string `json:"content"`
	IsReply   bool   `json:"isReply"`
}

func (c *CommentController) Post() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := postCommentReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.State = "error_req"
		return
	}
	if !(bson.IsObjectIdHex(req.ContentID) && (bson.IsObjectIdHex(req.FatherID) || req.FatherID == "")) {
		res.State = "error_id"
		return
	}
	if req.Content == "" {
		res.State = "null_content"
		return
	}
	err := c.Service.AddComment(c.Session.GetString("id"), req.ContentID, req.FatherID, req.Content, req.IsReply)
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

func (c *CommentController) DeleteBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = "error_id"
		return
	}
	err := c.Service.DeleteComment(id, c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}
