package controllers

import (
	"github.com/XMatrixStudio/Coffee/app/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/XMatrixStudio/Coffee/app/models"
)

type UsersController struct {
	Ctx     iris.Context
	Service services.UserService
	Session *sessions.Session
}

type commonRes struct {
	State string
	Data  string
}

func (c *UsersController) GetLogin() (res commonRes) {
	redirectURL := c.Ctx.FormValue("RedirectURL")
	url, state := c.Service.GetLoginURL(redirectURL)
	res.State = "success"
	res.Data = url
	c.Session.Set("state", state)
	return
}

type loginReq struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

func (c *UsersController) PostLogin() (res commonRes) {
	req := loginReq{}
	c.Ctx.ReadJSON(&req)
	if req.State != c.Session.GetString("state") {
		res.State = "error"
		res.Data = "State is error."
		return
	}
	userID, err := c.Service.LoginByCode(req.Code)
	if err != nil {
		res.State = "error"
		res.Data = err.Error()
		return
	}

	c.Session.Set("id", userID)
	res.State = "success"
	return
}

type userInfoRes struct {
	State        string
	Data		 string
	Email        string
	Name         string
	Class        int
	Info         models.UserInfo
	LikeNum      int64
	MaxSize      int64    // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize     int64    // 存储库已用大小 单位为KB
	SingleSize   int64    // 单个资源最大上限 -1为无上限
	FilesClass   []string // 文件分类
	ContentCount int64    // 内容数量
}

func (c *UsersController) GetInfo() (res userInfoRes) {
	if c.Session.Get("id") == nil {
		res.State = "error"
		res.Data = "not_login"
		return
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.State = "error"
		res.Data = "not_user"
		return
	}
	res.State = "success"
	res.Name = user.Name
	res.Info = user.Info
	res.Email = user.Email
	res.Class = user.Class
	res.LikeNum = user.LikeNum
	res.MaxSize = user.MaxSize
	res.UsedSize = user.UsedSize
	res.SingleSize = user.SingleSize
	res.FilesClass = user.FilesClass
	res.ContentCount = user.ContentCount
	return
}
