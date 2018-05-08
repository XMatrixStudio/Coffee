package content

// File 文件系统数据
type File struct {
	File  string `bson:"file"`  // File 文件系统的路径
	Size  int64  `bson:"size"`  // FileSize 文件大小
	Title string `bson:"title"` // Name 文件名（标题，并非真实文件米）
	Time  int64  `bson:"time"`  // Time 上传时间
	Count int64  `bson:"count"` // Count 下载次数
}
