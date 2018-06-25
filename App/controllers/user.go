package controllers

import (
	"regexp"

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

// GetLogin GET /user/login 获取登陆页面链接
func (c *UsersController) GetLogin() (res CommonRes) {
	redirectURL := c.Ctx.FormValue("RedirectURL")
	url, state := c.Service.GetLoginURL(redirectURL)
	res.State = "success"
	res.Data = url
	c.Session.Set("state", state)
	return
}

// loginReq 登陆请求
type loginReq struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// PostLogin POST /user/login 用户登陆
func (c *UsersController) PostLogin() (res CommonRes) {
	req := loginReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = "bad_request"
		return
	}
	if req.State != c.Session.GetString("state") {
		res.State = "error_state"
		return
	}
	userID, err := c.Service.LoginByCode(req.Code)
	if err != nil {
		res.State = err.Error()
		return
	}
	c.Session.Set("id", userID)
	res.State = "success"
	return
}

// UserInfoRes 用户信息返回
type UserInfoRes struct {
	State        string
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

// GetInfo GET /user/info 获取用户信息
func (c *UsersController) GetInfo() (res UserInfoRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	user, err := c.Service.GetUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.State = "not_user"
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

// PostInfo POST /user/info 更新用户信息
func (c *UsersController) PostInfo() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	err := c.Service.UpdateUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

type nameReq struct {
	Name string `json:"name"`
}

// PostName POST /user/name 更新用户名
func (c *UsersController) PostName() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := nameReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Name == "" || len(req.Name) > 20 {
		res.State = "error_name"
		return
	}
	if m, _ := regexp.MatchString(`[\\\/\(\)<|> "'{}:;]`, req.Name); m {
		res.State = "error_name"
		return
	}
	err = c.Service.UpdateUserName(c.Session.GetString("id"), req.Name)
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

// PostLogout POST /user/logout 退出登陆
func (c *UsersController) PostLogout() (res CommonRes) {
	c.Session.Delete("id")
	res.State = "success"
	return
}
