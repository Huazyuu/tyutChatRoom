package service

import "gin-gorilla/service/userService"

type ServiceGroup struct {
	UserService userService.UserService
}

var ServiceApp = new(ServiceGroup)
