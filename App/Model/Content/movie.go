package content

// Movie 电影
type Movie struct {
	File    File     `bson:"file"`    // File 文件系统
	URL     string   `bson:"url"`     // URL 在线地址
	Image   []string `bson:"Image"`   // 预览图
	Quality string   `bson:"quality"` // 画质
	Length  string   `bson:"length"`  // 时长
	Type    string   `bson:"type"`    // 类型
	Detail  string   `bson:"detail"`  // 介绍链接 (后期可以接入豆瓣API，自动获取电影详情)

}
