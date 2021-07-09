package app

import (
	"ginApp/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errorCode int, data interface{})  {
	g.C.JSON(httpCode, gin.H{
		"code": errorCode,
		"msg": e.GetMsg(errorCode),
		"data": data,
	})

	return
}
