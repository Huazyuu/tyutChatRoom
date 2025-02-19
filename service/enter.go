package service

import "gin-gorilla/service/userServer"

type ServiceGroup struct {
	UserService userServer.UserService
}

var ServiceApp = new(ServiceGroup)
