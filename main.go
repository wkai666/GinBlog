package main

import (
	"fmt"
	"ginApp/pkg/setting"
	"ginApp/routers"
	"net/http"
)

func main() {
	router := routers.InitRouter()
	//router.GET("/test", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "test",
	//	})
	//})

	s := &http.Server{
		Addr : fmt.Sprintf(":%d", setting.HTTPPort),
		Handler: router,
		ReadTimeout: setting.ReadTimeout,
		WriteTimeout: setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}