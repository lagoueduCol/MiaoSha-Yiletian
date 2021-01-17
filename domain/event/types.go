package event

import "github.com/letian0805/seckill/domain/product"

type Event struct {
	Topic  string    `json:"topic"`
	Banner string    `json:"banner"`
	List   []*Detail `json:"list"`
	Count  int       `json:"count"`
}

type Detail struct {
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	List      []*Goods `json:"list"`
	Count     int      `json:"count"`
}

type Goods struct {
	product.Goods
	EventPrice string `json:"event_price"`
	EventType  string `json:"event_type"`
}

type Info struct {
	EventPrice string `json:"event_price"`
	EventType  string `json:"event_type"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}
