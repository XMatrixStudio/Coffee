package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

// UsersController Users控制
type UsersController struct {
	Ctx     iris.Context
	Service services.UserService
	Session *sessions.Session
}

// GetLogin GET /login 获取登陆页面链接
func (c *UsersController) GetLogin() (res CommonRes) {
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

// PostLogin POST /login 用户登陆
func (c *UsersController) PostLogin() (res CommonRes) {
	req := loginReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = "error"
		res.Data = "Bad request."
		return
	}
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
	Data         string
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

// GetInfo GET /info 获取用户信息
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
	res.Name = user.Info.Name
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

func (c *UsersController) PostLogout() string {
	c.Session.Delete("id")
	return "success"
}
