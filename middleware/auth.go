package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"main/constants"
	"main/internal/jwt/private"
	"main/internal/user/service"
	"net/http"
	"strings"
	"sync"
)

var (
	syncOnce sync.Once
	svc      service.Interface
)

type AuthMiddleware struct {
	service.Interface
}

func NewAuthMiddleware(authService service.Interface) *AuthMiddleware {
	syncOnce.Do(func() {
		svc = authService
	})

	return &AuthMiddleware{
		svc,
	}
}

func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: missing Bearer token"})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := parseJWT(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		userDetails := claims.UserDetails
		if !a.IsUserValid(ctx, userDetails.UserID) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or inactive user"})
			ctx.Abort()
			return
		}

		ctx.Set(constants.PrivateUserDetails, &userDetails)
		ctx.Next()
	}
}

func parseJWT(tokenStr string) (*private.Claims, error) {
	claims := &private.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return viper.GetString("jwt.access_secret"), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
