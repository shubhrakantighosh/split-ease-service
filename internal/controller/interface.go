package controller

import (
	"github.com/gin-gonic/gin"
)

type Interface interface {
	LoginUser(ctx *gin.Context)
	RegisterUser(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)
	SendActivationEmail(ctx *gin.Context)
	ActivateUser(ctx *gin.Context)

	CreateGroup(ctx *gin.Context)
	UpdateGroup(ctx *gin.Context)
	RemoveGroup(ctx *gin.Context)
	GetUserGroups(ctx *gin.Context)
	//GetGroupDetails(ctx *gin.Context)
}
