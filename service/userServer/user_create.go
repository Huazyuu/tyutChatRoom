package userServer

import (
	"errors"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/utils"
	"gin-gorilla/utils/pwd"
)

func (UserService) CreateUser(username, password, email string) error {
	var userModel model.UserModel
	err := global.DB.Take(&userModel, "username= ? ", username).Error
	if err == nil {
		return errors.New("用户名已存在")
	}
	// hash pwd
	hashPwd, _ := pwd.HashPwd(password)

	// 生成图像并获取路径
	avatarPath, err := utils.DrawImage(username)
	if err != nil {
		return err
	}
	err = global.DB.Create(&model.UserModel{
		UserName: username,
		Password: hashPwd,
		Avatar:   avatarPath,
		Email:    email,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
