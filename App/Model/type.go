package model

// Album 相册
type Album struct {
	Images   []Image `bson:"images"`   // 图片预览图列表
	Title    string  `bson:"title"`    // 主题
	Time     int64   `bson:"time"`     // 时间
	Location string  `bson:"location"` // 地点
}

// App 应用
type App struct {
	File    File    `bson:"file"`    // 文件
	Web     string  `bson:"web"`     // 官方主页
	URL     string  `bson:"url"`     // 下载页面或地址
	Image   []Image `bson:"image"`   // 略缩图
	des     string  `bson:"des"`     // 使用说明
	Version string  `bson:"version"` // 本地当前版本
}

// File 文件系统数据
type File struct {
	File  string `bson:"file"`  // File 文件系统的路径
	Size  int64  `bson:"size"`  // FileSize 文件大小
	Title string `bson:"title"` // Name 文件名（标题，并非真实文件米）
	Time  int64  `bson:"time"`  // Time 上传时间
	Count int64  `bson:"count"` // Count 下载次数
	Type  string `bson:"type"`  // 文件类型
}

// Image 图片
type Image struct {
	Native bool   `bson:"native"` // 是否本地资源
	File   File   `bson:"file"`   // 文件地址
	URL    string `bson:"url"`    // 在线地址
	Thumb  string `bson:"thumb"`  // 缩略图
}

// Movie 电影
type Movie struct {
	File    File    `bson:"file"`    // 文件系统
	URL     string  `bson:"url"`     // 在线地址
	Image   []Image `bson:"Image"`   // 预览图
	Type    string  `bson:"type"`    // 类型
	Detail  string  `bson:"detail"`  // 介绍链接 (后期可以接入豆瓣API，自动获取电影详情)
	Watched bool    `bson:"watched"` // 是否已看（个人属性）
}
