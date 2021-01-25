package event

import "github.com/letian0805/seckill/domain/product"

type Topic struct {
	ID     string   `json:"id"`
	Topic  string   `json:"topic"`
	Banner string   `json:"banner"`
	List   []*Event `json:"list"`
	Count  int      `json:"count"`
}

type Event struct {
	ID        string   `json:"id"`
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
