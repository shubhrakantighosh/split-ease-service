package controller

import (
	"github.com/gin-gonic/gin"
	"main/constants"
	"main/internal/controller/adapter"
	"main/internal/controller/request"
	"main/internal/jwt/private"
	"main/util"
	"net/http"
	"strconv"
)

func (ctrl *Controller) CreateGroup(ctx *gin.Context) {
	var req request.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if svcErr := ctrl.groupService.CreateGroup(ctx, userID, req); svcErr.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": svcErr.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Group created successfully"})
}

func (ctrl *Controller) UpdateGroup(ctx *gin.Context) {
	var req request.UpdateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupID, parseErr := util.ParseUint(ctx.Param(constants.GroupID))
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	if svcErr := ctrl.groupService.UpdateGroup(ctx, userID, groupID, req); svcErr.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": svcErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group updated successfully"})
}

func (ctrl *Controller) RemoveGroup(ctx *gin.Context) {
	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupID, parseErr := util.ParseUint(ctx.Param(constants.GroupID))
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	if svcErr := ctrl.groupService.RemoveGroup(ctx, userID, groupID); svcErr.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": svcErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group removed successfully"})
}

func (ctrl *Controller) GetUserGroups(ctx *gin.Context) {
	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// later add filter
	groups, groupPermissions, err := ctrl.groupService.GetUserGroupsWithPermissions(ctx, userID, make(map[string]any))
	if err.Exists() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, adapter.BuildGroupPermissionsResponse(groups, groupPermissions))
}

func (ctrl *Controller) GetGroupDetails(ctx *gin.Context) {

}

func (ctrl *Controller) CreateGroupBill(ctx *gin.Context) {
	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupID, convErr := strconv.ParseUint(ctx.Param(constants.GroupID), 10, 64)
	if convErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req request.CreateBillRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	if err = ctrl.groupService.CreateGroupBill(ctx, userID, groupID, req); err.Exists() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Bill created successfully"})
}

func (ctrl *Controller) UpdateGroupBill(ctx *gin.Context) {
	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupID, convErr := strconv.ParseUint(ctx.Param(constants.GroupID), 10, 64)
	if convErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	billID, convErr := strconv.ParseUint(ctx.Param(constants.BillID), 10, 64)
	if convErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	var req request.UpdateBillRequest
	if bindErr := ctx.ShouldBindJSON(&req); bindErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	if err = ctrl.groupService.UpdateGroupBill(ctx, userID, groupID, billID, req); err.Exists() {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bill updated successfully"})
}

func (ctrl *Controller) DeleteGroupBill(ctx *gin.Context) {
	userID, err := private.GetUserID(ctx)
	if err.Exists() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	groupID, convErr := strconv.ParseUint(ctx.Param(constants.GroupID), 10, 64)
	if convErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	billID, convErr := strconv.ParseUint(ctx.Param(constants.BillID), 10, 64)
	if convErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	if err = ctrl.groupService.DeleteGroupBill(ctx, userID, groupID, billID); err.Exists() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bill deleted successfully"})
}
