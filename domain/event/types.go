package event

import "github.com/letian0805/seckill/domain/product"

type Topic struct {
	ID        int64    `json:"id"`
	Topic     string   `json:"topic"`
	Banner    string   `json:"banner"`
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	List      []*Event `json:"list"`
}

type Event struct {
	ID        int64    `json:"id"`
	StartTime int64    `json:"start_time"`
	EndTime   int64    `json:"end_time"`
	List      []*Goods `json:"list"`
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
