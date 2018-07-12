package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"jindan/datastruct"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	redisClient *redis.Client
	mgoSession  *mgo.Session
	mgoDataBase = "bcjys_test"
	getTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=ww338b579ec6e2589b&corpsecret=PunPJ7c-cvjAg_ew_JXPGWE18r_OfiGfwAFqjTqIjo0"
	weiXinUrl   = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
)

const (
	accessToken = "access_token"
	toparty     = "5"
	totag       = "1"
	agentid     = 1000002
	mgoUrl      = "mongodb://test_bc_guest:0IZ8v4s4pznzWCVz@39.107.159.28:3718/bcjys_test"
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "192.168.100.114:6379",
		Password: "redis_password", // no password set
		DB:       2,                // use default DB
	})
	getSession()
}

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mgoUrl)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

// 执行脚本
func GoOrderLoan() {

	log.Println("开始执行脚本")
	getAccessToken()
	ScriptDataAnalysis()
	log.Println("Sleep ")
	time.Sleep(time.Minute * 1)
}

// 执行脚本
func getAccessToken() {
	token, err := redisClient.Get(accessToken).Result()
	if err != nil {
		fmt.Println("err:", err)
	}
	if len(token) == 0 {
		response, err := http.Get(getTokenUrl)
		body, err := ioutil.ReadAll(response.Body) //[]uint8
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		tokenInfo := new(datastruct.RedisAccessToken)
		err = json.Unmarshal(body, tokenInfo)
		tokenValue := tokenInfo.Access_token
		err = redisClient.Set("access_token", tokenValue, 7200*time.Second).Err()
		if err != nil {
			panic(err)
		}
	}
}

func ScriptDataAnalysis() {
	mgoSession.SetMode(mgo.Monotonic, true)
	db := mgoSession.DB(mgoDataBase)
	apifinancingZJ0001Log := db.C("apifinancing_ZJ0001_log")
	timeStr := time.Now().Unix() - 86400
	fmt.Println("timeStr:", timeStr)
	iter := apifinancingZJ0001Log.Find(bson.M{"sign_status": 1, "is_update_mysql": bson.M{"$gte": 1}, "timeStr": bson.M{"$gte": timeStr}}).Iter()
	content := new(datastruct.MongoOrderloan)
	var successCount = 0
	var failCount = 0
	var failAssetsOrderNo string
	if iter != nil {
		for iter.Next(content) {
			Data := content.Data
			IsUpdateMysql := content.Is_update_mysql
			AssetsOrderNo := gjson.Get(Data, "assetOrderNo").String()
			fmt.Println("AssetsOrderNo:", AssetsOrderNo)
			if IsUpdateMysql == 1 {
				successCount++
			} else {
				failAssetsOrderNo = failAssetsOrderNo + "," + AssetsOrderNo
				failCount++
			}
		}
	}
	//fasong
	now := time.Now()
	formatNow := now.Format("2006-01-02 15:04:05")
	postContent := "时间：" + formatNow + "\n成功执行订单条数：" + strconv.Itoa(successCount) + "\n失败执行订单条数：" + strconv.Itoa(failCount) + "\n失败的订单号为：" + failAssetsOrderNo
	formt := `
    {
        "touser" : "",
        "toparty" : ` + toparty + `,
        "totag" : ` + totag + `,
        "msgtype" : "text",
        "agentid" : ` + strconv.Itoa(agentid) + `,
        "text" : {
            "content" : "%s"。
        },
        "safe":0
    }`
	postBody := fmt.Sprintf(formt, postContent)
	jsonValue := []byte(postBody)
	tokenValue, err := redisClient.Get(accessToken).Result()
	fmt.Println("tokenValue:", tokenValue)
	resp, err := http.Post(weiXinUrl+tokenValue, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	fmt.Println("resp:", resp)
	iter.Close()
}
