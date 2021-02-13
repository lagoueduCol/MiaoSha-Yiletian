package stock

import (
	"reflect"
	"testing"

	"github.com/letian0805/seckill/infrastructure/stores/redis"
)

func TestStock(t *testing.T) {
	var (
		st  Stock
		err error
		val int64
	)
	if err = redis.Init(); err != nil {
		t.Fatal(err)
	}
	if st, err = NewRedisStock("101", "1001"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cli := redis.GetClient()
		cli.Del("seckill#101#1001")
		cli.Del("seckill#101#1001#123")
	}()
	if err = st.Set(10, 100); err != nil {
		t.Fatal(err)
	}
	if val, err = st.Get(); err != nil {
		t.Fatal(err)
	} else if val != 10 {
		t.Fatal("not equal 10")
	}
	if val, err = st.Sub("123"); err != nil {
		t.Fatal(err)
	} else if val != 9 {
		t.Fatal("not equal 9", val)
	}
	if err = st.Del(); err != nil {
		t.Fatal(err)
	}
	if val, err = st.Get(); err != nil {
		t.Fatal(err)
	} else if val != 0 {
		t.Fatal("not equal 0")
	}
}

func TestRedis(t *testing.T) {
	redis.Init()
	cli := redis.GetClient()
	script := `
		return redis.call('get', 'seckill#101#1001')
	`
	res, err := cli.Eval(script, []string{}).Result()
	t.Log(res, reflect.TypeOf(res), err)
}
