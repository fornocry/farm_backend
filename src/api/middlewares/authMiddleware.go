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

func (m MiddlewareServiceImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			pkg.PanicException(constant.Unauthorized, "Invalid authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := pkg.VerifyJwtToken(tokenString)
		if err != nil {
			pkg.PanicException(constant.Unauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			pkg.PanicException(constant.Unauthorized, "Invalid token claims")
			return
		}

		userAuthIDStr, ok := claims["user_auth_id"].(string)
		if !ok {
			pkg.PanicException(constant.Unauthorized, "Invalid user_auth_id in token")
			return
		}

		userAuthID, err := uuid.Parse(userAuthIDStr)
		if err != nil {
			pkg.PanicException(constant.Unauthorized, "Invalid user_auth_id format")
			return
		}

		userAuth, err := m.userRepository.GetByAuthId(userAuthID)
		if err != nil {
			pkg.PanicException(constant.Unauthorized, "User not found")
			return
		}

		c.Set("user", userAuth.User)
		c.Set("user_auth", userAuth)
		c.Next()
	}
}

func MiddlewareServiceInit(userRepository repository.UserRepository) *MiddlewareServiceImpl {
	return &MiddlewareServiceImpl{
		userRepository: userRepository,
	}
}
