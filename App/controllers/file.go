package controllers

import (
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
)

// FileController 上传接口
type FileController struct {
	Ctx     iris.Context
	Service services.FileService
	Session *sessions.Session
}

// PostFileMeta 添加文件元信息
func (c *FileController) PostMeta() (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	req := services.UploadMeta{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.State = StatusBadReq
		return
	}
	// 实现秒传
	if c.Service.IsExist(req.MD5) {
		if err := c.Service.AddOwn(req.MD5, c.Session.GetString("id")); err != nil {
			res.State = StatusNotAllow
			return
		}
		res.State = StatusExist
		return
	}
	if err := c.Service.AddFile(req); err != nil {
		res.State = StatusNotAllow
		return
	}
	res.State = StatusSuccess
	return
}

func (c *FileController) PostUpload() {
	if c.Session.Get("id") == nil {
		c.Ctx.StatusCode(iris.StatusForbidden)
		return
	}
	file, _, err := c.Ctx.FormFile("file")
	if err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	token := c.Ctx.FormValue("token")
	index, err := strconv.Atoi(c.Ctx.FormValue("index"))
	if len(token) != 32 || err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := c.Service.UploadFileToTemp(&file, token, index); err != nil {
		c.Ctx.StatusCode(iris.StatusBadRequest)
		c.Ctx.WriteString(err.Error())
		return
	}
	return
}

func (c *FileController) PostMergeBy(token string) (res CommonRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	if len(token) != 32 {
		res.State = StatusBadReq
		return
	}
	if err := c.Service.MergeFile(token, c.Session.GetString("id")); err != nil {
		res.State = StatusNotAllow
		res.Data = err.Error()
		return
	}
	res.State = StatusSuccess
	return
}

type FileMetaRes struct {
	State string `json:"status"`
	Data  services.UploadMeta
}

func (c *FileController) GetMetaBy(token string) (res FileMetaRes) {
	if c.Session.Get("id") == nil {
		res.State = StatusNotLogin
		return
	}
	if len(token) != 32 {
		res.State = StatusBadReq
		return
	}
	if data, err := c.Service.GetInfo(token); err != nil {
		res.State = StatusNotAllow
		return
	} else {
		res.Data = data
	}
	res.State = StatusSuccess
	return
}
