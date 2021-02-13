package stock

import (
	"errors"

	"github.com/letian0805/seckill/infrastructure/stores"
)

type memStock struct {
	eventID string
	goodsID string
	key     string
}

var (
	cache = stores.NewIntCache()

	ErrNotFound = errors.New("not found")
)

func NewMemStock(eventID string, goodsID string) (Stock, error) {
	if eventID == "" || goodsID == "" {
		return nil, errors.New("invalid event id or goods id")
	}
	key := eventID + "." + goodsID
	stock := &memStock{
		eventID: eventID,
		goodsID: goodsID,
		key:     key,
	}
	return stock, nil
}

func (ms *memStock) Set(val int64, expiration int64) error {
	cache.Set(ms.key, val)
	return nil
}

func (ms *memStock) Get() (int64, error) {
	val, ok := cache.Get(ms.key)
	if !ok {
		return 0, ErrNotFound
	}
	return val, nil
}

func (ms *memStock) Sub(uid string) (int64, error) {
	return cache.Add(ms.key, -1), nil
}

func (ms *memStock) Del() error {
	cache.Del(ms.key)
	return nil
}

func (ms *memStock) EventID() string {
	return ms.eventID
}

func (ms *memStock) GoodsID() string {
	return ms.goodsID
}
