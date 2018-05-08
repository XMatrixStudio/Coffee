package content

// App 应用
type App struct {
	File    File    `bson:"file"`    // 文件
	Web     string  `bson:"web"`     // 官方主页
	URL     string  `bson:"url"`     // 下载页面或地址
	Image   []Image `bson:"image"`   // 略缩图
	des     string  `bson:"des"`     // 使用说明
	Version string  `bson:"version"` // 本地当前版本
}
