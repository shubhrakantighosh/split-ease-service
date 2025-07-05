package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	authService "main/internal/auth/service"
	"main/middleware"
	opostgres "main/pkg/db/postgres"
)

func PublicRoutes(ctx context.Context, s *gin.Engine) {
	route := s.Group("api/v1")

	route.POST("/login")

	authMiddleware := middleware.NewAuthMiddleware(authService.Wire(ctx, opostgres.GetCluster().DbCluster))
	baseRoute := s.Group("/api/v1", middleware.SanitizeQueryParams(), authMiddleware.Authenticate())

	fmt.Println(baseRoute)

}
