package main

import (
	"errors"
	"goPipeline/stream"
	"goPipeline/utils"
	"goPipeline/web"
	"os"
	"runtime"

	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/app"
)

func main() {
	if runtime.GOOS == "windows" {
		os.Setenv("ZONEINFO", "zoneinfo.zip")
		realPort := utils.FindPort("0.0.0.0", utils.MakeRange(10000, 20000))
		utils.RegisterPort(web.WebServerPort, realPort)
		web.WebServerPort = realPort
	}
	var comp stream.Component
	args := os.Args
	switch args[1] {
	case "DataConnector":
		// 数据连接器
		comp = &stream.DataConnectorComponent{}
	default:
		errors.New("not support")
	}
	comp.InitHandler()
	comp.SioHandler()
	app.Run(comp.CallHandler)
}
