package content

// Movie 电影
type Movie struct {
	File    File    `bson:"file"`    // 文件系统
	URL     string  `bson:"url"`     // 在线地址
	Image   []Image `bson:"Image"`   // 预览图
	Type    string  `bson:"type"`    // 类型
	Detail  string  `bson:"detail"`  // 介绍链接 (后期可以接入豆瓣API，自动获取电影详情)
	Watched bool    `bson:"watched"` // 是否已看（个人属性）
}
