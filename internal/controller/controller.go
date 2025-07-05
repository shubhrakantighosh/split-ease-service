package controller

import (
	"github.com/gin-gonic/gin"
	"main/internal/auth/service"
)

type Controller struct {
	service.Interface
}

func NewController(s service.Interface) *Controller {
	return &Controller{
		s,
	}
}

func (ctrl *Controller) Ge(ctx *gin.Context) {
	ctrl.Get(ctx)
}
