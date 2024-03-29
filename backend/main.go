package main

import (
	"bufio"
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
	cmdPython := exec.Command("python3", "scripts/pyRuntime.py")

	pyStdout, _ := cmdPython.StdoutPipe()

	go func() {
		scanner := bufio.NewScanner(pyStdout)
		for scanner.Scan() {
			m := scanner.Text()
			log.Infof("Python脚本编辑器消息：%s", m)
		}
	}()

	pyStderr, _ := cmdPython.StderrPipe()

	go func() {
		scanner := bufio.NewScanner(pyStderr)
		for scanner.Scan() {
			m := scanner.Text()
			log.Errorf("Python脚本编辑器报错：%s", m)
		}
	}()
	errPython := cmdPython.Start()
	if errPython != nil {
		log.Errorf("启动fastapi失败，失败原因：%s", errPython.Error())
	}
	if val, ok := config.GetArgs()["--__python__pkgs"]; ok {
		if len(val) > 0 {
			for _, pkg := range strings.Split(val, ",") {
				cmd := exec.Command("pip", "install", pkg, "-i", "https://mirrors.aliyun.com/pypi/simple")
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

	appMode := "edit"
	if val, ok := config.GetArgs()["--__app_mode"]; ok {
		appMode = val
	}
	var comp stream.Component
	args := os.Args
	if len(args) >= 2 {
		// 数据连接器 DataConnector
		if utils.SlicesContain([]string{"DataConnector"}, args[1]) {
			comp = &stream.DataConnectorComponent{Type: args[1], Mode: appMode}
		} else {
			panic("启动组件名错误")
		}
	} else {
		panic("未提供启动组件名称")
	}
	comp.InitHandler()
	comp.SioHandler()
	app.Run(comp.CallHandler)
}
