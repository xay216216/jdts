package datastruct

import "time"

// mysql -----------------------------------
// 订单
type BcOrders struct {
	OrderId            int       `orm:"not null pk autoincr INT(10)" grom:"primary_key"`
	OrderNo            string    `orm:"VARCHAR(32)"`
	AssetsOrderNo      string    `orm:"VARCHAR(64)"`
	LoanStatus         int       `orm:"not null default 0 index TINYINT(2)"`
	AssetName          string    `orm:"VARCHAR(64)"`
	FinancingOrderNo   string    `orm:"VARCHAR(64)"`
	FinancingId        string    `orm:"CHAR(6)"`
	FinancingName      string    `orm:"VARCHAR(64)"`
	ProductName        string    `orm:"VARCHAR(64)"`
	AssetsId           string    `orm:"CHAR(6)"`
	ProductId          string    `orm:"CHAR(9)"`
	BorrowMoney        int       `orm:"INT(10)"`
	OrderStatus        int       `orm:"not null default 0 index TINYINT(2)"`
	RepaymentMethod    int       `orm:"not null default 0 index TINYINT(2)"`
	LoanMoney          int       `orm:"INT(10)"`
	LoanTime           time.Time `orm:"TIMESTAMP"`
	InterestStart      time.Time `orm:"TIMESTAMP"`
	InterestEnd        time.Time `orm:"TIMESTAMP"`
	BorrowDeadline     int       `orm:"not null default 0 index TINYINT(2)"`
	Periods            int       `orm:"not null default 1 TINYINT(3)"`
	DeadlineUnit       int       `orm:"not null default 0 index TINYINT(2)"`
	InterestRate       int       `orm:"INT(10)"`
	InterestRateType   int       `orm:"not null default 0 index TINYINT(2)"`
	TotalPrincipal     int       `orm:"INT(10)"`
	TotalInterests     int       `orm:"INT(10)"`
	TotalServiceCharge int       `orm:"INT(10)"`
	TotalMoney         int       `orm:"INT(10)"`
	PackageId          string    `orm:"INT(10)"`
	Realname           string    `orm:"VARCHAR(64)"`
}

// 放款结果
type BcOrderResults struct {
	Cid      int    `orm:"not null pk default 0 INT(10)"`
	Src      string `orm:"default '' VARCHAR(100)"`
	Video    string `orm:"default '' VARCHAR(200)"`
	VideoImg string `orm:"default '' VARCHAR(200)"`
	Content  string `orm:"TEXT"`
	Stat     int    `orm:"not null default 1 TINYINT(3)"`
	Did      int    `orm:"INT(10)"`
	Author   string `orm:"default '' VARCHAR(300)"`
}

// 订单状态日志
type BcOrderStatusLogs struct {
	Id       int       `orm:"not null pk autoincr INT(10)"`
	ChanId   int       `orm:"not null default 0 index INT(10)"`
	Cid      int       `orm:"not null default 0 index INT(10)"`
	Top      int       `orm:"not null default 0 INT(10)"`
	Pub      int       `orm:"not null default 0 TINYINT(3)"`
	PubTime  time.Time `orm:"TIMESTAMP"`
	Stat     int       `orm:"not null default 1 TINYINT(3)"`
	Order    int       `orm:"not null default 0 INT(10)"`
	CardType int       `orm:"not null default 0 INT(10)"`
}

//出借人信息
type BcAssetsOrdersLenders struct {
	Cid      int    `orm:"not null pk default 0 INT(10)"`
	Image    string `orm:"default '' VARCHAR(255)"`
	Url      string `orm:"not null default '' VARCHAR(300)"`
	Duration string `orm:"not null default '' VARCHAR(20)"`
	Stat     int    `orm:"not null default 1 TINYINT(3)"`
	Did      int    `orm:"default 0 INT(10)"`
}

//出借人信息
type BcFinancingsAssetsRelations struct {
	Cid      int    `orm:"not null pk default 0 INT(10)"`
	Image    string `orm:"default '' VARCHAR(255)"`
	Url      string `orm:"not null default '' VARCHAR(300)"`
	Duration string `orm:"not null default '' VARCHAR(20)"`
	Stat     int    `orm:"not null default 1 TINYINT(3)"`
	Did      int    `orm:"default 0 INT(10)"`
}
