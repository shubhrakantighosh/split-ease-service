package controller

import (
	"github.com/gin-gonic/gin"
)

type Interface interface {
	LoginUser(ctx *gin.Context)
	RegisterUser(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	SendActivationEmail(ctx *gin.Context)
	ActivateUser(ctx *gin.Context)

	CreateGroup(ctx *gin.Context)
	UpdateGroup(ctx *gin.Context)
	RemoveGroup(ctx *gin.Context)
	GetUserGroups(ctx *gin.Context)
	GetGroupDetails(ctx *gin.Context)

	CreateGroupBillForUser(ctx *gin.Context)
	AssignUserToGroup(ctx *gin.Context)
	UpdateGroupBill(ctx *gin.Context)
	DeleteGroupBill(ctx *gin.Context)

	CalculateBillSplits(ctx *gin.Context)
	RecalculateBillSplits(ctx *gin.Context)
}
