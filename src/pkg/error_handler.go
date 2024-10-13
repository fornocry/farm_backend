package pkg

import (
	"crazyfarmbackend/src/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func PanicHandler(c *gin.Context) {
	if err := recover(); err != nil {
		str := fmt.Sprint(err)
		strArr := strings.SplitN(str, ":", 2)
		key := strArr[0]
		message := strArr[1]
		switch key {
		case constant.DataNotFound.GetResponseStatus():
			c.JSON(http.StatusBadRequest, BuildResponse_(key, message))
			c.Abort()
		case constant.Unauthorized.GetResponseStatus():
			c.JSON(http.StatusUnauthorized, BuildResponse_(key, message))
			c.Abort()
		case constant.WrongDataBody.GetResponseStatus():
			c.JSON(http.StatusBadRequest, BuildResponse_(key, message))
			c.Abort()
		case constant.WrongBody.GetResponseStatus():
			c.JSON(http.StatusBadRequest, BuildResponse_(key, message))
			c.Abort()
		default:
			c.JSON(http.StatusInternalServerError, BuildResponse_(key, message))
			c.Abort()
		}
	}
}
