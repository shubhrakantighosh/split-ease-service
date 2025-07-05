package router

import (
	"context"
	"github.com/gin-gonic/gin"
	ctrl "main/internal/controller"
	userService "main/internal/user/service"
	"main/middleware"
	opostgres "main/pkg/db/postgres"
)

func RegisterPublicRoutes(ctx context.Context, engine *gin.Engine) {
	apiV1 := engine.Group("/api/v1", middleware.RequestLogger())

	userController := ctrl.Wire(ctx, opostgres.GetCluster().DbCluster)

	// Public user routes
	userRoutes := apiV1.Group("/users")
	{
		userRoutes.POST("/register", userController.RegisterUser)
		userRoutes.POST("/login", userController.LoginUser)
		userRoutes.POST("/activate", userController.ActivateUser)
		userRoutes.POST("/send-activation", userController.SendActivationEmail)
	}

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(userService.Wire(ctx, opostgres.GetCluster().DbCluster))

	// Protected routes
	protectedRoutes := apiV1.Group("/", middleware.SanitizeQueryParams(), authMiddleware.Authenticate())
	{
		protectedRoutes.PUT("/users", userController.UpdateUserProfile)
	}

	groupRoutes := apiV1.Group("/groups", middleware.SanitizeQueryParams(), authMiddleware.Authenticate())
	{
		// Group management
		groupRoutes.POST("/", userController.CreateGroup)
		groupRoutes.PUT("/:group_id", userController.UpdateGroup)
		groupRoutes.DELETE("/:group_id", userController.RemoveGroup)
		groupRoutes.GET("/", userController.GetUserGroups)

		// Bills for group
		groupRoutes.POST("/:group_id/bills", userController.CreateGroupBill)
		groupRoutes.PUT("/:group_id/bills/:bill_id", userController.UpdateGroupBill)
		groupRoutes.DELETE("/:group_id/bills/:bill_id", userController.DeleteGroupBill)

		// Permissions (optional if needed to expose)
		// groupRoutes.GET("/:group_id/permissions", userController.GetGroupPermissions)
		// groupRoutes.POST("/:group_id/permissions", userController.UpdateGroupPermissions)
	}
}
