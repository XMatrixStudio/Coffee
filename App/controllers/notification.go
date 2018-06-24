package controllers

import (
	"github.com/kataras/iris"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris/sessions"
	"strconv"
	"github.com/XMatrixStudio/Coffee/App/models"
)

// NotificationController Like
type NotificationController struct {
	Ctx     iris.Context
	Service services.NotificationService
	Session *sessions.Session
}

type readReq struct {
	IsRead bool `json:"isRead"`
}

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

type NotificationRes struct {
	State string
	Notification []models.NotificationDetail
}

func (c *NotificationController) GetAll() (res NotificationRes) {
	if c.Session.Get("id") == nil {
		res.State = "not_login"
		return
	}
	var err error
	res.Notification, err = c.Service.GetUserNotification(c.Session.GetString("id"))
	if err != nil {
		res.State = err.Error()
		return
	}
	res.State = "success"
	return
}