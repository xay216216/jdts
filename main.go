package main

import (
	//"encoding/json"
	//"fmt"
	//"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	//_ "github.com/go-sql-driver/mysql"
	//"jindan/datastruct"
	//"jindan/order"
	//"log"
)

func init() {
	// set default database
	//orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8", 30)
	// register model
	//orm.RegisterModel(new(datastruct.BcOrders))
	// create table
	//orm.RunSyncdb("default", false, true)
	//mongo

}

func main() {
	//创建定时任务
	scriptthreeNewBot := toolbox.NewTask("scriptthree", "0/5 * * * * *", newsBot)
	//添加定时任务
	toolbox.AddTask("scriptthree", scriptthreeNewBot)
	//启动定时任务
	toolbox.StartTask()
	defer toolbox.StopTask()
	select {}
}

func newsBot() error {
	order.UpdateOrderStatus()
	return nil
}
