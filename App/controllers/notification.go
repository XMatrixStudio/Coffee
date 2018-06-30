package controllers

import (
	"strconv"

	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

// NotificationController Like
type NotificationController struct {
	Ctx     iris.Context
	Service services.NotificationService
	Session *sessions.Session
}

// readReq 标记已读请求
type readReq struct {
	IsRead bool `json:"isRead"`
}

// PatchReadBy PATCH /notification/read/{NotificationID} 标记指定通知为已读
func (c *NotificationController) PatchReadBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	req := readReq{}
	err := c.Ctx.ReadJSON(&req)
	if err != nil {
		res.State = "err_req"
		return
	}
	err = c.Service.SetRead(c.Session.GetString("id"), id, req.IsRead)
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

// DeleteBy DELETE /notificaiton/{NotificationID} 删除指定通知
func (c *NotificationController) DeleteBy(id string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	err := c.Service.RemoveByID(c.Session.GetString("id"), id)
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}

// GetUnread GET /notification/unerad 获取未读通知数
func (c *NotificationController) GetUnread() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	count, err := c.Service.GetUserUnreadCount(c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	res.Data = strconv.Itoa(count)
	return
}

// NotificationRes 通知集合返回值
type NotificationRes struct {
	State        string
	Notification []services.NotificationData
}

// GetAll GET /notification/all 获取用户所有通知
func (c *NotificationController) GetAll() (res NotificationRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	res.Notification, _ = c.Service.GetUserNotification(c.Session.GetString("id"))
	res.State = "success"
	return
}
