package datastruct

// mongo -----------------------------------
// 订单
type RedisAccessToken struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}
