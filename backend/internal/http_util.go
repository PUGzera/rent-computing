package http_util

import (
	"errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetUserIdFromCtx(ctx *gin.Context, jwtHandler jwt.GinJWTMiddleware) (string, error) {
	claims, err := jwtHandler.GetClaimsFromJWT(ctx)
	if err != nil {
		return "", err
	}

	idClaim, ok := claims["id"]
	if !ok {
		return "", errors.New("claim \"id\" not found")
	}

	id, ok := idClaim.(string)
	if !ok {
		return "", errors.New("id was not of type string")
	}

	return id, nil
}
