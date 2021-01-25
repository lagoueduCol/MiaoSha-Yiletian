package stock

import (
	"errors"
	"fmt"
	"time"

	"github.com/letian0805/seckill/infrastructure/stores/redis"
)

type Stock interface {
	// 设置库存，并设置过期时间
	Set(val int64, expire int64) error
	// 直接返回剩余库存
	Get() (int64, error)
	// 尝试扣减一个库存，并返回剩余库存
	Sub() (int64, error)
	// 删除库存数据
	Del() error
	// 返回活动 ID
	EventID() string
	// 返回商品 ID
	GoodsID() string
}

type redisStock struct {
	eventID string
	goodsID string
	key     string
}

func NewRedisStock(eventID string, goodsID string) (Stock, error) {
	if eventID == "" || goodsID == "" {
		return nil, errors.New("invalid event id or goods id")
	}
	stock := &redisStock{
		eventID: eventID,
		goodsID: goodsID,
		key:     fmt.Sprintf("seckill#%s#%s", eventID, goodsID),
	}

	return stock, nil
}

func (rs *redisStock) Set(val int64, expiration int64) error {
	cli := redis.GetClient()
	return cli.Set(rs.key, val, time.Duration(expiration)*time.Second).Err()
}

func (rs *redisStock) Sub() (int64, error) {
	cli := redis.GetClient()
	return cli.Decr(rs.key).Result()
}

func (rs *redisStock) Get() (int64, error) {
	cli := redis.GetClient()
	if val, err := cli.Get(rs.key).Int64(); err != nil && err != redis.Nil {
		return 0, err
	} else {
		return val, nil
	}
}

func (rs *redisStock) Del() error {
	cli := redis.GetClient()
	return cli.Del(rs.key).Err()
}

func (rs *redisStock) EventID() string {
	return rs.eventID
}

func (rs *redisStock) GoodsID() string {
	return rs.goodsID
}
