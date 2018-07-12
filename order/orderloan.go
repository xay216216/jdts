package order

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"jindan/datastruct"
	"jindan/util"
	"log"
	"time"
)

// 移动数量
var moveCount = 0

func init() {
	//mongo
	sessionMongo, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer sessionMongo.Close()
	sessionMongo.SetMode(mgo.Monotonic, true)
	db := sessionMongo.DB("test")
	//mysql
	maxIdle := 30 // 参数4(可选)  设置最大空闲连接
	maxConn := 30 // 参数5(可选)  设置最大数据库连接 (go >= 1.2)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8", maxIdle, maxConn)
	if err != nil {
		panic(err)
	}
	defer sessionMysql.Close()
}

// 执行脚本
func UpdateOrderStatus() {
	defer func() {
		if err := recover(); err != nil {
			util.DingDingNotice("markdown", "脚本任务任务崩溃了", "脚本任务任务崩溃了")
		}
	}()

	for {
		log.Println("开始执行脚本")

		findMoveData()

		log.Println("Sleep ")
		time.Sleep(time.Minute * 5)
	}
}

// 查询需要移动的数据
func findMoveData() {
	apifinancingZJ0001Log := db.C("apifinancing_ZJ0001_log")
	iter := apifinancingZJ0001Log.Find(bson.M{"is_update_mysql": 0, "sign_status": 1}).Limit(200).Iter()
	content := new(datastruct.MomgoOrderloan)
	if iter != nil {
		for iter.Next(content) {
			Data := content.Data
			Id := content.Id
			fmt.Println("Name:", Id)
			AssetsOrderNo := gjson.Get(Data, "assetOrderNo")
			PaymentStatus := gjson.Get(Data, "paymentStatus")
			Sid := gjson.Get(Data, "sid")
			PlanList := gjson.Get(Data, "planList")
			//fmt.Println("Name:", PlanList)
			ProductId = getProductId(Sid)
			if len(ProductId) == 0 {
				apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 5, "update_time": time.Now()}})
				continue
			}
			if len(PlanList) == 0 {
				apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 4, "update_time": time.Now()}})
				continue
			}
			orm.RegisterModel(new(datastruct.BcOrders), new(datastruct.BcOrderResults), new(datastruct.BcOrderStatusLogs), new(datastruct.BcAssetsOrdersLenders), new(datastruct.BcFinancingsAssetsRelations))
			o := orm.NewOrm()
			OrderInfo := BcOrders{AssetsOrderNo: AssetsOrderNo, ProductId: ProductId, OrderStatus: 60}
			err := o.Read(&OrderInfo)
			if err == orm.ErrNoRows {
				apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 3, "update_time": time.Now()}})
				fmt.Println("未找到此债权ID对应的订单!:", AssetsOrderNo)
				continue
			} else if err == orm.ErrMissPK {
				apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 3, "update_time": time.Now()}})
				fmt.Println("主键未找到此债权ID对应的订单!:", AssetsOrderNo)
				continue
			} else {
				if 1 == PaymentStatus {
					Result = addOrderResults(content, OrderInfo)
					if Result {
						apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 1, "update_time": time.Now(), "product_id": ProductId}})
						continue
					} else {
						apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 6, "update_time": time.Now()}})
						fmt.Println("写入订单结果表mysql操作失败!:", AssetsOrderNo)
						continue
					}
				} else {
					OrderInfo.OrderStatus = 71
					if num, err := o.Update(&OrderInfo); err == nil {

						apifinancingZJ0001Log.UpdateId(content.Id, bson.M{"$set": bson.M{"is_update_mysql": 2, "update_time": time.Now()}})
						fmt.Println("未找到此债权ID对应的订单!:", AssetsOrderNo)
						continue
					}
				}
			}
		}
	}
}

func addOrderResults(content, OrderInfo map[string]interface{}) int {
	// 组合插入数据
	check := 1 - channel.IsCheck
	params := map[string]interface{}{
		"title":         content.Title,
		"desc":          content.Summary,
		"check":         check,
		"src":           content.Source,
		"content":       content.Content,
		"author":        content.Author,
		"channels":      channel.AppChannel,
		"content_style": channel.ContentCss,
		"mongo_id":      content.Id.Hex(),
		"thumb":         content.Thumb,
		"card_type":     channel.CardType,
	}
}

func getProductId(x int) string {
	var (
		sidToProduct = map[int]string{
			59: "ZC0006-01",
			50: "ZC0016-01",
			83: "ZC0118-01",
			55: "ZC0119-01",
			82: "ZC0122-01",
			88: "ZC0123-01",
			56: "ZC0125-01",
			81: "ZC0126-01",
			72: "ZC0127-01",
			76: "ZC0128-01",
			73: "ZC0129-01",
			58: "ZC0130-01",
			84: "ZC0131-01",
		}
	)
	for k, v := range sidToProduct {
		if k == x {
			return v
		}
	}
	return ""
}
