package services

import (
	"github.com/XMatrixStudio/Coffee/app/models"
	"github.com/XMatrixStudio/Violet.SDK.Go"
)

type UserService interface {
	InitViolet(c violetSdk.Config)
	// 登陆部分API
	/*GetLoginUrl() string
	Login(name, password string) (valid bool, email string, err error)
	GetUser(code string) (ID, name string, err error)
	Register(name, email, password string) (err error)
	GetEmailCode(email string) error
	ValidEmail(email, vCode string) error
	GetUserInfo(id string) (user models.Users, err error)
	SetUserName(id, name string) error*/
}

type userService struct {
	Model  models.UserModel
	Violet violetSdk.Violet
}

func (s *userService) InitViolet(c violetSdk.Config) {
	s.Violet = violetSdk.NewViolet(c)
}
