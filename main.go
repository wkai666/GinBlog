package main

import (
	"fmt"
	"ginApp/models"
	"ginApp/pkg/gredis"
	"ginApp/pkg/logging"
	"ginApp/pkg/setting"
	"ginApp/routers"
	"github.com/fvbock/endless"
	"log"
	"syscall"
)

func main() {

	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeOut
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeOut
	endless.DefaultHammerTime = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func (add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server error: %v", err)
	}
}