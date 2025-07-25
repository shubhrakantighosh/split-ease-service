package service

import (
	billSvc "main/internal/bill/service"
	groupRepo "main/internal/group/repository"
	groupPermissionSvc "main/internal/group_permission/service"
	userSvc "main/internal/user/service"
	"sync"
)

type Service struct {
	groupRepo          groupRepo.Interface
	groupPermissionSvc groupPermissionSvc.Interface
	billSvc            billSvc.Interface
	userSvc            userSvc.Interface
}

var (
	syncOnce sync.Once
	svc      *Service
)

func NewService(
	groupRepo groupRepo.Interface,
	groupPermissionSvc groupPermissionSvc.Interface,
	billSvc billSvc.Interface,
	userSvc userSvc.Interface,
) *Service {
	syncOnce.Do(func() {
		svc = &Service{groupRepo: groupRepo, groupPermissionSvc: groupPermissionSvc, billSvc: billSvc, userSvc: userSvc}
	})

	return svc
}
