package content

// Album 相册
type Album struct {
	Images   []Image `bson:"images"`   // 图片预览图列表
	Title    string  `bson:"title"`    // 主题
	Time     int64   `bson:"time"`     // 时间
	Location string  `bson:"location"` // 地点
}
