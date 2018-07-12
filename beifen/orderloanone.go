package script

import (
	"fmt"
	//"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"jindan/datastruct"
	"jindan/util"
	"log"
	"time"
)

/*func init() {
    //mongo
    sessionMongo, err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }
    defer sessionMongo.Close()
    sessionMongo.SetMode(mgo.Monotonic, true)
}*/

// 执行脚本
func GoOrderLoan() {
	defer func() {
		if err := recover(); err != nil {
			util.DingDingNotice("markdown", "脚本任务任务崩溃了", "脚本任务任务崩溃了")
		}
	}()

	for {
		log.Println("开始执行脚本")

		ScriptDataAnalysis()

		log.Println("Sleep ")
		time.Sleep(time.Minute * 1)
	}
}

func ScriptDataAnalysis() {
	//mongo
	sessionMongo, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer sessionMongo.Close()
	sessionMongo.SetMode(mgo.Monotonic, true)
	db := sessionMongo.DB("test")
	apifinancingZJ0001Log := db.C("apifinancing_ZJ0001_log")
	timeStr := time.Now().Format("2006-01-02")
	now := time.Now().UTC()
	// 显示时间格式： UnixDate = "Mon Jan _2 15:04:05 MST 2006"
	fmt.Printf("%s\n", now.Format(time.UnixDate))
	// 显示时间戳
	fmt.Printf("%ld\n", now.Unix())
	// 显示时分:Kitchen = "3:04PM"
	fmt.Printf("%s\n", now.Format("3:04PM"))
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 23:59:59", time.Local)
	afterPubTime, _ := time.Parse("2006-01-02", timeStr)
	fmt.Println("start_time:", afterPubTime)
	fmt.Println("start_time:", startTime)
	fmt.Println("end_time:", endTime)
	//iter := apifinancingZJ0001Log.Find(bson.M{"is_update_mysql": 9, "sign_status": 1}).Iter()
	//iter := apifinancingZJ0001Log.Find(bson.M{"sign_status": 1, "update_time": bson.M{"$gte": startTime, "$lte": endTime}}).Iter()
	iter := apifinancingZJ0001Log.Find(bson.M{"sign_status": 1, "update_time": bson.M{"$gte": startTime}}).Iter()
	content := new(datastruct.MomgoOrderloan)
	if iter != nil {
		for iter.Next(content) {
			Data := content.Data
			//Id := content.Id
			fmt.Println("Name:", Data)
		}
	}
}
