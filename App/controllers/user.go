package controllers

import (
	"regexp"

	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"html/template"
	"github.com/globalsign/mgo/bson"
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
	ID string
	State        string
	Email        string
	Name         string
	Class        int
	Info         models.UserInfo
	MaxSize      int64    // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize     int64    // 存储库已用大小 单位为KB
	SingleSize   int64    // 单个资源最大上限 -1为无上限
}

func (c *UsersController) GetInfoBy(id string) (res UserInfoRes) {
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
	user, err := c.Service.GetUserInfo(id)
	if err != nil {
		res.State = "not_user"
		return
	}
	res.State = "success"
	res.ID = user.ID.Hex()
	res.Name = user.Info.Name
	res.Info = user.Info
	res.Email = user.Email
	res.Class = user.Class
	if isOwn {
		res.MaxSize = user.MaxSize
		res.UsedSize = user.UsedSize
		res.SingleSize = user.SingleSize
	}
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
	req.Name = template.HTMLEscapeString(req.Name)
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
