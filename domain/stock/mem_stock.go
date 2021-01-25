package stock

import (
	"errors"
	"sync/atomic"

	"github.com/letian0805/seckill/infrastructure/utils"
)

type memStock struct {
	val     int64
	eventID string
	goodsID string
	key     string
}

var cache = utils.NewMemCache(10)

func NewMemStock(eventID string, goodsID string) (Stock, error) {
	if eventID == "" || goodsID == "" {
		return nil, errors.New("invalid event id or goods id")
	}
	key := eventID + "." + goodsID
	stock, ok := cache.Get(key).(*memStock)
	if !ok {
		stock = &memStock{
			eventID: eventID,
			goodsID: goodsID,
			key:     key,
		}
		cache.Set(stock.key, stock)
	}
	return stock, nil
}

func (ms *memStock) Set(val int64, expiration int64) error {
	atomic.StoreInt64(&ms.val, val)
	return cache.Expire(ms.key, expiration)
}

func (ms *memStock) Get() (int64, error) {
	return atomic.LoadInt64(&ms.val), nil
}

func (ms *memStock) Sub() (int64, error) {
	return atomic.AddInt64(&ms.val, -1), nil
}

func (ms *memStock) Del() error {
	atomic.StoreInt64(&ms.val, 0)
	cache.Del(ms.key)
	return nil
}

func (ms *memStock) EventID() string {
	return ms.eventID
}

func (ms *memStock) GoodsID() string {
	return ms.goodsID
}
