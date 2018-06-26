package coffee

import (
	"github.com/kataras/iris"

	"time"

	"github.com/XMatrixStudio/Coffee/App/controllers"
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/XMatrixStudio/Coffee/App/services"
	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)


// Config 配置文件
type Config struct {
	Mongo  models.Mongo     `yaml:"Mongo"`  // mongoDB配置
	Server ServerConfig     `yaml:"Server"` // iris配置
	Violet violetSdk.Config `yaml:"Violet"` // Violet配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"Host"` // 服务器监听地址
	Port string `yaml:"Port"` // 服务器监听端口
	Dev  bool   `yaml:"Dev"`  // 是否开发环境
}

// RunServer 开始运行服务
func RunServer(c Config) {
	// 初始化数据库
	Model, err := models.NewModel(c.Mongo)
	if err != nil {
		panic(err)
	}
	// 初始化服务
	// 启动服务器
	app := iris.New()
	if c.Server.Dev {
		app.Logger().SetLevel("debug")
	}

	sessionManager := sessions.New(sessions.Config{
		Cookie:  "sessionCoffee",
		Expires: 24 * time.Hour,
	})

	Service := services.NewService(Model)

	users := mvc.New(app.Party("/user"))
	userService := Service.GetUserService()
	userService.InitViolet(c.Violet)
	users.Register(userService, sessionManager.Start)
	users.Handle(new(controllers.UsersController))

	content := mvc.New(app.Party("/content"))
	content.Register(Service.GetContentService(), sessionManager.Start)
	content.Handle(new(controllers.ContentController))

	comment := mvc.New(app.Party("/comment"))
	comment.Register(Service.GetCommentService(), sessionManager.Start)
	comment.Handle(new(controllers.CommentController))

	like := mvc.New(app.Party("/like"))
	like.Register(Service.GetLikeService(), sessionManager.Start)
	like.Handle(new(controllers.LikeController))

	notification := mvc.New(app.Party("/notification"))
	notification.Register(Service.GetNotificationService(), sessionManager.Start)
	notification.Handle(new(controllers.NotificationController))

	app.Run(
		// Starts the web server
		iris.Addr(c.Server.Host+":"+c.Server.Port),
		// Disables the updater.
		iris.WithoutVersionChecker,
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
