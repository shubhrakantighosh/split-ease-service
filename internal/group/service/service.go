package service

import (
	billSvc "main/internal/bill/service"
	groupRepo "main/internal/group/repository"
	groupPermissionSvc "main/internal/group_permission/service"
	"sync"
)

type Service struct {
	groupRepo          groupRepo.Interface
	groupPermissionSvc groupPermissionSvc.Interface
	billSvc            billSvc.Interface
}

var (
	syncOnce sync.Once
	svc      *Service
)

func NewService(
	groupRepo groupRepo.Interface,
	groupPermissionSvc groupPermissionSvc.Interface,
	billSvc billSvc.Interface,
) *Service {
	syncOnce.Do(func() {
		svc = &Service{groupRepo: groupRepo, groupPermissionSvc: groupPermissionSvc, billSvc: billSvc}
	})

	return svc
}
