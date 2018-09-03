package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

// LikeController Like
type LikeController struct {
	Ctx     iris.Context
	Service services.LikeService
	Session *sessions.Session
}

// LikeRes 用户点赞数据返回值
type LikeRes struct {
	State string
	Data  []string
}

// Get GET /like 获取用户点赞列表
func (c *LikeController) Get() (res LikeRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	var err error
	res.Data, err = c.Service.GetUserLikes(c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
	}
	res.State = StatusSuccess
	return
}

type likeReq struct {
	IsContent bool `json:"isContent"`
	IsComment bool `json:"isComment"`
	IsReply   bool `json:"isReply"`
}

// PostBy POST /like/{contentID} 对某个内容点赞
func (c *LikeController) PostBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = StatusBadReq
		return
	}
	req := likeReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = StatusBadReq
		return
	}
	if req.IsContent {
		err = c.Service.AddLikeToContent(id, c.Session.GetString("id"))
	} else if req.IsComment {
		err = c.Service.AddLikeToComment(id, c.Session.GetString("id"), false)
	} else if req.IsReply {
		err = c.Service.AddLikeToComment(id, c.Session.GetString("id"), true)
	} else {
		res.State = StatusBadReq
		return
	}
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = StatusSuccess
	return
}

// PatchBy PATCH /like/{contentID} 取消用户对某个内容的点赞
func (c *LikeController) PatchBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	if !bson.IsObjectIdHex(id) {
		res.State = StatusBadReq
		return
	}
	req := likeReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = StatusBadReq
		return
	}
	if req.IsContent {
		err = c.Service.RemoveLikeFromContent(id, c.Session.GetString("id"))
	} else if req.IsComment {
		err = c.Service.RemoveLikeFromComment(id, c.Session.GetString("id"), false)
	} else if req.IsReply {
		err = c.Service.RemoveLikeFromComment(id, c.Session.GetString("id"), true)
	} else {
		res.State = StatusBadReq
		return
	}
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = StatusSuccess
	return
}
