package shop

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/letian0805/seckill/infrastructure/pool"

	"github.com/letian0805/seckill/domain/stock"
	"github.com/letian0805/seckill/infrastructure/utils"

	"github.com/sirupsen/logrus"

	"github.com/letian0805/seckill/infrastructure/mq"
)

const (
	OK         = 0
	ErrNoStock = 1001
	ErrRedis   = 1002
	ErrTimeout = 1003

	requestTimeout = 60
)

type Context struct {
	Request *http.Request
	Conn    net.Conn
	Writer  *bufio.ReadWriter
	GoodsID string
	EventID string
	UID     string
}

var queue mq.Queue

func Init() {
	queueFactory := mq.NewFactory("memory")
	if queueFactory == nil {
		panic("no memory queue factory")
	}
	queue, _ = queueFactory.New("shop")
	go func() {
		for {
			task, err := queue.Consume()
			if err != nil {
				logrus.Error(err)
				break
			}
			task.Do()
		}
	}()
}

func Handle(ctx *Context) {
	start := time.Now().Unix()
	t := func() {
		data := &utils.Response{
			Code: OK,
			Data: nil,
			Msg:  "ok",
		}
		status := http.StatusOK
		now := time.Now().Unix()
		if now-start > requestTimeout {
			data.Msg = "request timeout"
			data.Code = ErrTimeout
		} else {
			// 扣减 Redis 库存
			st, _ := stock.NewRedisStock(ctx.EventID, ctx.GoodsID)
			if s, err := st.Sub(ctx.UID); err != nil {
				data.Msg = err.Error()
				data.Code = ErrRedis
			} else if s < 0 {
				data.Msg = "no stock"
				data.Code = ErrNoStock
			}
		}
		// 此处实现操作购物车的逻辑

		body, _ := json.Marshal(data)
		resp := &http.Response{
			Proto:         ctx.Request.Proto,
			ProtoMinor:    ctx.Request.ProtoMinor,
			ProtoMajor:    ctx.Request.ProtoMajor,
			Header:        make(http.Header),
			ContentLength: int64(len(body)),
			Body:          ioutil.NopCloser(bytes.NewReader(body)),
			StatusCode:    status,
			Close:         false,
		}
		resp.Header.Set("Content-Type", "application/json")
		resp.Write(ctx.Writer)
		ctx.Writer.Flush()
		ctx.Conn.Close()
	}

	queue.Produce(pool.TaskFunc(t))
}
