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
	Sub(uid string) (int64, error)
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

func (rs *redisStock) Sub(uid string) (int64, error) {
	cli := redis.GetClient()
	script := `
	local history=redis.call('get',KEYS[1])
	local stock=redis.call('get', KEYS[2])
	if (history and history >= '1') or stock==false or stock <= '0' then
		return -1
	else
		stock=redis.call('decr', KEYS[2])
		if stock >= 0 and redis.call('set', KEYS[1], '1', 'ex', 86400) then
			return stock
		else
			return -1
		end
	end`
	if res, err := cli.Eval(script, []string{fmt.Sprintf("%s#%s", rs.key, uid), rs.key}).Result(); err != nil {
		return -1, err
	} else if resInt, ok := res.(int64); ok && resInt != -1 {
		return resInt, nil
	} else {
		return -1, errors.New("redis error")
	}
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
