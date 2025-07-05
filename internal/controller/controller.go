package controller

import (
	"github.com/gin-gonic/gin"
	userService "main/internal/user/service"
)

type Controller struct {
	userSvc userService.Interface
}

func NewController(userSvc userService.Interface) *Controller {
	return &Controller{
		userSvc: userSvc,
	}
}

func (ctrl *Controller) Login(ctx *gin.Context) {
	ctrl.userS

}
