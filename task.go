package main

import (
	"github.com/astaxie/beego/toolbox"
	"jdts/conf"
	"jdts/script"
)

func main() {
	//创建定时任务
	myConfig := new(conf.Config)
	myConfig.InitConfig("./conf/app.conf")
	taskTime := myConfig.Read("task", "taskTime")
	scriptthreeNewBot := toolbox.NewTask("scriptthree", taskTime, newsBot)
	//添加定时任务
	toolbox.AddTask("scriptthree", scriptthreeNewBot)
	//启动定时任务
	toolbox.StartTask()
	defer toolbox.StopTask()
	select {}
}

func newsBot() error {
	script.GoOrderLoan()
	return nil
}
