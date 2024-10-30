package model

type Trade struct {
	Id             string `json:"id"`             // 成交单id
	TradePair      string `json:"tradePair"`      // 交易对
	MakerId        string `json:"makerId"`        // maker订单id
	TakerId        string `json:"takerId"`        // taker订单id
	MakerUser      int64  `json:"makerUser"`      // maker用户id
	TakerUser      int64  `json:"takerUser"`      // taker用户id
	Price          string `json:"price"`          // 成交价
	Quantity       string `json:"quantity"`       // 成交数量
	TakerOrderSide string `json:"takerOrderSide"` // taker订单方向
	TakerOrderType string `json:"takerOrderType"` // taker订单类型
	Timestamp      int64  `json:"timestamp"`      // 成交时间
}
