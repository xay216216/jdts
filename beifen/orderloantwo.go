package script

import (
	"bytes"
	//"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"jindan/datastruct"
	"jindan/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	redisClient *redis.Client
	getTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ww338b579ec6e2589b&corpsecret=PunPJ7c-cvjAg_ew_JXPGWE18r_OfiGfwAFqjTqIjo0"
	weiXinUrl   = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
)

const (
	accessToken = "access_token"
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})

}

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

// 执行脚本
func getAccessToken() {
	token, err := redisClient.Get(accessToken).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("token type:%T\n", token)
	fmt.Println("sss:", token)
	if len(token) == 0 {
		fmt.Println("body2:", 222)
		response, err := http.Get(getTokenUrl)
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		//fmt.Println("body:", body)
		defer response.Body.Close()
		tokenValue := gjson.Get(body, "access_token").String()
		//tokenInfo := new(datastruct.RedisAccessToken)
		//err = json.Unmarshal(body, tokenInfo)
		//tokenValue := tokenInfo.Access_token
		fmt.Println("accessTokenssssss:", tokenValue)
		fmt.Printf("v1 type:%T\n", tokenValue)
		err = redisClient.Set("access_token", tokenValue, 7200*time.Second).Err()
		if err != nil {
			panic(err)
		}
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
	timeStr := time.Now().Unix() - 86400
	//log.Println("Id:", timeStr)
	iter := apifinancingZJ0001Log.Find(bson.M{"sign_status": 1, "is_update_mysql": bson.M{"$gte": 1}, "timeStr": bson.M{"$gte": timeStr}}).Iter()
	content := new(datastruct.MomgoOrderloan)

	var successCount = 0
	var failCount = 0
	//var sucessData []string
	//var failData []string //切片
	var failAssetsOrderNo string
	if iter != nil {
		for iter.Next(content) {
			Data := content.Data
			Id := content.Id
			IsUpdateMysql := content.Is_update_mysql
			AssetsOrderNo := gjson.Get(Data, "assetOrderNo").String()
			log.Println("Id:", Id)
			log.Println("IsUpdateMysql:", IsUpdateMysql)
			fmt.Println("Name:", AssetsOrderNo)
			if IsUpdateMysql == 1 {
				successCount++
			} else {
				//failData = append(failData, AssetsOrderNo)
				failAssetsOrderNo = failAssetsOrderNo + "," + AssetsOrderNo
				failCount++
			}
		}
	}
	//fasong
	postContent := "成功执行订单条数：" + strconv.Itoa(successCount) + "\n失败执行订单条数：" + strconv.Itoa(failCount) + "\n失败的订单号为：" + failAssetsOrderNo
	formt := `
    {
        "touser" : "XiaoAYong",
        "toparty" : "3",
        "totag" : "1",
        "msgtype" : "text",
        "agentid" : 1000002,
        "text" : {
            "content" : "%s"。
        },
        "safe":0
    }`
	postBody := fmt.Sprintf(formt, postContent)
	jsonValue := []byte(postBody)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	tokenValue, err := redisClient.Get(accessToken).Result()
	fmt.Println("tokenValuessss:", tokenValue)
	resp, err := http.Post(weiXinUrl+tokenValue, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	fmt.Println("resp:", resp)
	iter.Close()
	defer redisClient.Close()
}
