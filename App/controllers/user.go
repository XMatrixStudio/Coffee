package controllers

import (
	"regexp"

	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/v12"
	"html/template"
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
	ID         string
	State      string
	Email      string
	Name       string
	Class      int
	Info       models.UserInfo
	MaxSize    int64 // 存储库使用最大上限 -1为无上限 单位为KB
	UsedSize   int64 // 存储库已用大小 单位为KB
	SingleSize int64 // 单个资源最大上限 -1为无上限
}

func (c *UsersController) GetInfoBy(id string) (res UserInfoRes) {
	isOwn := false
	if id == "self" {
		if c.Session.Get("id") == nil {
			res.State = StatusNotLogin
			return
		}
		id = c.Session.GetString("id")
		isOwn = true
	} else if !bson.IsObjectIdHex(id) {
		res.State = StatusBadReq
		return
	}
	user, err := c.Service.GetUserInfo(id)
	if err != nil {
		res.State = StatusNotAllow
		return
	}
	res.State = StatusSuccess
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
		res.State = StatusNotLogin
		return
	}
	err := c.Service.UpdateUserInfo(c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = StatusSuccess
	return
}

type nameReq struct {
	Name string `json:"name"`
}

// PostName POST /user/name 更新用户名
func (c *UsersController) PostName() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	req := nameReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Name == "" || len(req.Name) > 20 {
		res.State = StatusBadReq
		return
	}
	if m, _ := regexp.MatchString(`[\\\/\(\)<|> "'{}:;]`, req.Name); m {
		res.State = StatusBadReq
		return
	}
	req.Name = template.HTMLEscapeString(req.Name)
	err = c.Service.UpdateUserName(c.Session.GetString("id"), req.Name)
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = StatusSuccess
	return
}

// PostLogout POST /user/logout 退出登陆
func (c *UsersController) PostLogout() (res CommonRes) {
	c.Session.Delete("id")
	res.State = StatusSuccess
	return
}

// LoginReq POST /user/login 登陆请求
type LoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// PostLogin POST /user/login/pass 密码登陆
func (c *UsersController) PostLoginPass() (result CommonRes) {
	req := LoginReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil || req.Name == "" || req.Password == "" {
		result.State = StatusBadReq
		return
	}
	valid, data, err := c.Service.Login(req.Name, req.Password)
	if err != nil { // 与Violet连接发生错误
		result.State = StatusError
		result.Data = err.Error()
		return
	}
	if !valid { // 用户邮箱未激活
		c.Session.Set("email", data)
		result.State = StatusNotValid
		result.Data = data
		return
	} else {
		c.Session.Delete("email")
	}

	userID, tErr := c.Service.LoginByCode(data)
	if tErr != nil { // 无法获取用户详情
		result.State = StatusError
		result.Data = tErr.Error()
		return
	}
	c.Session.Set("id", userID)
	result.State = StatusSuccess
	return
}

// RegisterReq POST /user/register 注册请求
type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PostRegister POST /user/register 注册
func (c *UsersController) PostRegister() (res CommonRes) {
	req := RegisterReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.State = StatusBadReq
	}
	if err := c.Service.Register(req.Name, req.Email, req.Password); err != nil {
		res.State = err.Error()
	} else {
		c.Session.Set("email", req.Name)
		res.State = StatusSuccess
	}
	return
}

// PostEmail POST /user/email 获取邮箱验证码
func (c *UsersController) PostEmail() (res CommonRes) {
	if c.Session.Get("email") == nil {
		res.State = StatusNotLogin
		return
	}
	if err := c.Service.GetEmailCode(c.Session.GetString("email")); err != nil {
		res.State = err.Error()
	} else {
		res.State = StatusSuccess
	}
	return
}

// ValidReq POST /user/valid/ 请求
type ValidReq struct {
	VCode string `json:"vCode"`
}

// PostValid POST /user/valid/ 验证邮箱
func (c *UsersController) PostValid() (res CommonRes) {
	if c.Session.Get("email") == nil {
		res.State = StatusNotLogin
		return
	}
	req := ValidReq{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.State = StatusBadReq
	}
	if err := c.Service.ValidEmail(c.Session.GetString("email"), req.VCode); err != nil {
		res.State = err.Error()
	} else {
		res.State = StatusSuccess
	}
	return
}
