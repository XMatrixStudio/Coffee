package controllers

import (
	"github.com/globalsign/mgo/bson"
	"github.com/XMatrixStudio/Coffee/App/services"
	"strings"
	"fmt"
)

func (c *ContentController) PostAlbum() (res CommonRes) {
	c.Ctx.SetMaxRequestBodySize(1024 * 1024 * 1024) // 1G
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := services.NewAlbumData{}
	if err := c.Ctx.ReadForm(&req); err != nil {
		res.State = err.Error()
		return
	}
	if err := c.Service.AddAlbum(c.Ctx, c.Session.GetString("id"), req); err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

func (c *ContentController) GetAlbumBy(id string) (res ContentsRes) {
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
	res.Data = c.Service.GetAlbumByUser(id, !isOwn)
	return
}

func (c *ContentController) GetFileBy(id string, filePath string) {
	if c.Session.Get("id") == nil {
		c.Ctx.WriteString("not_login")
		return
	}
	if !bson.IsObjectIdHex(id) {
		c.Ctx.WriteString("err_id")
		return
	}
	path := strings.Replace(filePath, "|", "/", -1)
	fmt.Println(path)
	name, err := c.Service.GetFile(c.Session.GetString("id"), id,  path)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	c.Ctx.SendFile(filePath, name)
	return
}