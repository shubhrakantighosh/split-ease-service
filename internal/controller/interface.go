package controller

import "github.com/gin-gonic/gin"

type Interface interface {
	Ge(ctx *gin.Context)
}
