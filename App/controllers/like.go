package controllers

import (
	"github.com/kataras/iris"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris/sessions"
	"github.com/globalsign/mgo/bson"
)

// LikeController Like
type LikeController struct {
	Ctx     iris.Context
	Service services.LikeService
	Session *sessions.Session
}


type LikeRes struct {
	State string
	Data []string
}

func (c *LikeController) Get() (res LikeRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	var err error
	res.Data, err = c.Service.GetUserLikes(c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
	}
	res.State = "success"
	return
}

type likeReq struct {
	IsContent bool `json:"isContent"`
	IsComment bool `json:"isComment"`
	IsReply  bool `json:"isReply"`
}

func (c *LikeController) PostBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = "error_id"
		return
	}
	req := likeReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = "error_req"
		return
	}
	if req.IsContent {
		err = c.Service.AddLikeToContent(id, c.Session.GetString("id"))
	} else if req.IsComment {
		err = c.Service.AddLikeToComment(id, c.Session.GetString("id"), false)
	} else if req.IsReply {
		err = c.Service.AddLikeToComment(id, c.Session.GetString("id"), true)
	} else {
		res.State = "error_req"
		return
	}
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

func (c *LikeController) PutBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = "error_id"
		return
	}
	req := likeReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = "error_req"
		return
	}
	if req.IsContent {
		err = c.Service.RemoveLikeFromContent(id, c.Session.GetString("id"))
	} else if req.IsComment {
		err = c.Service.RemoveLikeFromComment(id, c.Session.GetString("id"), false)
	} else if req.IsReply {
		err = c.Service.RemoveLikeFromComment(id, c.Session.GetString("id"), true)
	} else {
		res.State = "error_req"
		return
	}
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}
