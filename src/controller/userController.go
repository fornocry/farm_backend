package controller

import (
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController interface {
	AuthUser(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	GetMyUpgrades(c *gin.Context)
	GetMyFields(c *gin.Context)
	GetMyReferrals(c *gin.Context)
}

type UserControllerImpl struct {
	userService service.UserService
}

func (u UserControllerImpl) AuthUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse, err := u.userService.AuthUser(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, userResponse)
	return
}

func (u UserControllerImpl) GetMe(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse := u.userService.GetMe(c)
	c.JSON(http.StatusOK, userResponse)
	return
}
func (u UserControllerImpl) GetMyUpgrades(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse := u.userService.GetUserUpgrade(c)
	c.JSON(http.StatusOK, userResponse)
	return
}
func (u UserControllerImpl) GetMyFields(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse := u.userService.GetMyFields(c)
	c.JSON(http.StatusOK, userResponse)
	return
}
func (u UserControllerImpl) GetMyReferrals(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse := u.userService.GetMyReferrals(c)
	c.JSON(http.StatusOK, userResponse)
	return
}

func UserControllerInit(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userService: userService,
	}
}
