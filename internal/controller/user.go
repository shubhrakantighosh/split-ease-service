package controller

import (
	"github.com/gin-gonic/gin"
	"main/internal/controller/adapter"
	"main/internal/controller/request"
	"main/internal/jwt/private"
	"main/internal/model"
	"net/http"
)

func (ctrl *Controller) LoginUser(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.userSvc.AuthenticateUser(ctx, req.Email, req.Password)
	if err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, adapter.BuildAuthTokenResponse(token))
}

func (ctrl *Controller) RegisterUser(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userSvc.CreateUserAccount(ctx, req); err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (ctrl *Controller) UpdateUserProfile(ctx *gin.Context) {
	var req request.UpdateRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err = ctrl.userSvc.UpdateUserProfile(ctx, user); err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (ctrl *Controller) SendActivationEmail(ctx *gin.Context) {
	var req request.SendOTPRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userSvc.SendActivationEmail(ctx, req.Email); err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activation email sent"})
}

func (ctrl *Controller) ActivateUser(ctx *gin.Context) {
	var req request.ActivateRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userSvc.ActivateUserAccount(ctx, req.Email, req.Password, req.Otp); err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User activated successfully"})
}
