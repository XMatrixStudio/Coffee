package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/globalsign/mgo/bson"
)

func (c *ContentController) PostMovie() (res CommonRes) {
	c.Ctx.SetMaxRequestBodySize(4* 1024 * 1024 * 1024) // 4G
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := services.ContentData{}
	if err := c.Ctx.ReadForm(&req); err != nil {
		res.State = err.Error()
		return
	}
	if err := c.Service.AddMovie(c.Ctx, c.Session.GetString("id"), req); err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

func (c *ContentController) GetMovieBy(id string) (res ContentsRes) {
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
	res.Data = c.Service.GetMovieByUser(id, !isOwn)
	return
}
