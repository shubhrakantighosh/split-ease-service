package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main/constants"
	"main/internal/auth/service"
	"main/internal/jwt/private"
	"net/http"
	"strings"
	"sync"
	"time"
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
		claims, err := ParseJWT(tokenStr)
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

var jwtKey = []byte("your_secret_key") // Use env variable in production

func GenerateJWT(userID uint64) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenStr string) (*private.Claims, error) {
	claims := &private.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
