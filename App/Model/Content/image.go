package content

// Image 图片
type Image struct {
	File  File   `bson:"file"`  // 文件地址
	Thumb string `bson:"thumb"` // 缩略图
}
