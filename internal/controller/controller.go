package controller

import (
	billSplitSvc "main/internal/bill_split/service"
	groupService "main/internal/group/service"
	userService "main/internal/user/service"
	"sync"
)

type Controller struct {
	userSvc      userService.Interface
	groupService groupService.Interface
	billSplitSvc billSplitSvc.Interface
}

var (
	syncOnce sync.Once
	ctrl     *Controller
)

func NewController(
	userSvc userService.Interface,
	groupService groupService.Interface,
	billSplitSvc billSplitSvc.Interface,
) *Controller {
	syncOnce.Do(func() {
		ctrl = &Controller{
			userSvc:      userSvc,
			groupService: groupService,
			billSplitSvc: billSplitSvc,
		}
	})

	return ctrl
}
