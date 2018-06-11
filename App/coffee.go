package coffee

import (
	"github.com/kataras/iris"

	model "github.com/XMatrixStudio/Coffee/App/model"
)

// Config 配置文件
type Config struct {
	Mongo  model.Mongo  // mongoDB配置
	Server ServerConfig // iris配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string // 服务器监听地址
	Port string // 服务器监听端口
	Dev  bool   // 是否开发环境
}

// RunServer 开始运行服务
func RunServer(c Config) {
	model.InitMongo(c.Mongo) // 初始化数据库
	app := iris.New()
	if c.Server.Dev {
		app.Logger().SetLevel("debug")
	}
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
