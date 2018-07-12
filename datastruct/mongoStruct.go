package datastruct

//import "time"
import "gopkg.in/mgo.v2/bson"

// mongo -----------------------------------
// 订单
type MongoOrderloan struct {
	Id              bson.ObjectId `bson:"_id"`
	Financing_id    string        `json:"financing_id"`
	Time            string        `json:"time"`
	Time_str        int           `json:"timeStr"`
	Is_update_mysql int           `json:"is_update_mysql"`
	Sign            string        `json:"sign"`
	Sign_status     int           `json:"sign_status"`
	Data            string        `json:"data"`
	Product_id      string        `json:"product_id"`
}
