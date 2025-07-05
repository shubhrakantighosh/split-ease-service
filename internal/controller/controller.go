package controller

import (
	groupService "main/internal/group/service"
	userService "main/internal/user/service"
	"sync"
)

type Controller struct {
	userSvc      userService.Interface
	groupService groupService.Interface
}

var (
	syncOnce sync.Once
	ctrl     *Controller
)

func NewController(userSvc userService.Interface, groupService groupService.Interface) *Controller {
	syncOnce.Do(func() {
		ctrl = &Controller{
			userSvc:      userSvc,
			groupService: groupService,
		}
	})

	return ctrl
}
