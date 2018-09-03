package controllers

import (
	"fmt"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"io"
	"os"
)

// FileController 上传接口
type FileController struct {
	Ctx     iris.Context
	Service services.FileService
	Session *sessions.Session
}

// PostFileMeta 添加文件元信息
func (c *FileController) PostMeta() (res CommonRes) {
	req := services.UploadMeta{}
	if err := c.Ctx.ReadJSON(&req); err != nil {
		res.State = StatusBadReq
		return
	}

	return
}

func (c *FileController) PostUpload() {
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	md5 := c.Ctx.FormValue("token")
	fmt.Println(md5)
	index := c.Ctx.FormValue("index")
	fmt.Println(index)
	defer file.Close()

	// Create a file with the same name
	// assuming that you have a folder named 'uploads'
	out, err := os.OpenFile("./UserData/"+md5,
		os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	defer out.Close()

	io.Copy(out, file)
}
