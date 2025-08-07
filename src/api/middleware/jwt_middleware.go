package middleware

import (
	domain "bookem-user-service/domain"
	utils "bookem-user-service/util"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	ID       uint
	Username string
	Role     domain.UserRole
}

func GetJwtString(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return "", domain.ErrUnauthenticated
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("invalid authorization model (must be Bearer)")
	}

	jwt := strings.SplitN(header, "Bearer ", 2)[1]
	_, err := utils.ParseJWT(jwt)

	return jwt, err
}

func GetJwtData(ctx *gin.Context) (jwt.MapClaims, error) {
	jwtString, err := GetJwtString(ctx)
	if err != nil {
		return nil, err
	}

	jwt, err := utils.ParseJWT(jwtString)
	if err != nil {
		return nil, err
	}

	return jwt, err
}

// GetJwt returns the JWT data embedded in the request header. If the user is
// unauthenticated (no JWT in the request), the funciton returns (nil,
// ErrUnauthenticated).
func GetJwt(ctx *gin.Context) (*Jwt, error) {
	jwtData, err := GetJwtData(ctx)
	if err != nil {
		return nil, err
	}

	jwt := Jwt{
		ID:       jwtData["sub"].(uint),
		Username: jwtData["username"].(string),
		Role:     jwtData["role"].(domain.UserRole),
	}

	return &jwt, nil
}
