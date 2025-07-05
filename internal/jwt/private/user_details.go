package private

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"main/constants"
	"main/pkg/apperror"
	"net/http"
)

type Claims struct {
	UserDetails UserDetails `json:"user_details"`
	jwt.RegisteredClaims
}

type UserDetails struct {
	UserID uint64 `json:"user_id"`
}

func GetUserID(ctx *gin.Context) (uint64, apperror.Error) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return 0, apperror.NewWithMessage("user details missing in context", http.StatusUnauthorized)
	}

	return userDetails.UserID, apperror.Error{}
}
