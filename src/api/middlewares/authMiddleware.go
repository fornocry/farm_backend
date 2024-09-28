package middlewares

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
)

type MiddlewareService interface {
	AuthMiddleware() gin.HandlerFunc
	RequestIdMiddleware() gin.HandlerFunc
	CorsMiddleware() gin.HandlerFunc
}

type MiddlewareServiceImpl struct {
	userRepository repository.UserRepository
}

func (u MiddlewareServiceImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqToken := c.Request.Header.Get("Authorization")
		jwtSecretToken := strings.Split(reqToken, "Bearer ")[1]
		token, err := pkg.VerifyJwtToken(jwtSecretToken)
		if err != nil {
			pkg.PanicException(constant.Unauthorized, "")
		}
		userAuthId := token.Claims.(jwt.MapClaims)["user_auth_id"].(string)
		userAuthIdUuid, err := uuid.Parse(userAuthId)
		userAuth, err := u.userRepository.GetByAuthId(userAuthIdUuid)
		if err != nil {
			pkg.PanicException(constant.Unauthorized, "")
		}
		user := userAuth.User
		c.Set("user", user)
		c.Set("user_auth", userAuth)
		c.Next()
	}
}

func MiddlewareServiceInit(userRepository repository.UserRepository) *MiddlewareServiceImpl {
	return &MiddlewareServiceImpl{
		userRepository: userRepository,
	}
}
