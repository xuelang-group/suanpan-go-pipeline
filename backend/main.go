package main

import (
	"errors"
	"goPipeline/stream"
	"goPipeline/utils"
	"goPipeline/web"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/xuelang-group/suanpan-go-sdk/config"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/app"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/log"
)

func main() {
	if val, ok := config.GetArgs()["--__python__pkgs"]; ok {
		if len(val) > 0 {
			for _, pkg := range strings.Split(val, ",") {
				cmd := exec.Command("pip", "install", pkg, "-i", "https://pypi.mirrors.ustc.edu.cn/simple")
				log.Infof("开始安装python依赖库%s...", pkg)
				err := cmd.Run()
				if err != nil {
					log.Errorf("安装python依赖库%s失败，失败原因：%s", pkg, err.Error())
				} else {
					log.Infof("安装python依赖库%s成功", pkg)
				}
			}
		}
	}
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
