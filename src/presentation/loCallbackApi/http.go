package loCallbackApi

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func RegisterLoCallbackApi(ginRouterGroup *gin.RouterGroup) {
	ginRouterGroup.POST("/locallback", ReceiveLoCallback())
}

func ReceiveLoCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(c.Request.Body)
		str := buf.String()
		logrus.Infof("LoCallback: %s", str)
		c.JSON(http.StatusOK, "Callback received")
	}
}
