package services

// UserService 用户服务
type UserService interface {
	// GetLoginUrl 获取登陆跳转连接
	GetLoginUrl() string
	// Login 用户登录
	Login(token string) bool
}
