package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/mvcc/mvccpb"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/letian0805/seckill/infrastructure/stores/etcd"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel  string `json:"logLevel"`
	RateLimit struct {
		Middle int `json:"middle"`
		Low    int `json:"low"`
	} `json:"rateLimit"`
	CircuitBreaker struct {
		Cpu     int `json:"cpu"`
		Latency int `json:"latency"`
	} `json:"circuitBreaker"`
}

var configLock = &sync.RWMutex{}
var config = &Config{}

func WatchClusterConfig() error {
	cli := etcd.GetClient()
	key := "/seckill/config"
	resp, err := cli.Get(context.Background(), key)
	if err != nil {
		return err
	}
	update := func(kv *mvccpb.KeyValue) (bool, error) {
		if string(kv.Key) == key {
			var tmpConfig *Config
			err = json.Unmarshal(kv.Value, &tmpConfig)
			if err != nil {
				logrus.Error("update cluster config failed, error:", err)
				return false, err
			}
			configLock.Lock()
			*config = *tmpConfig
			logrus.Info("update cluster config ", *config)
			configLock.Unlock()
			return true, nil
		}
		return false, nil
	}
	for _, kv := range resp.Kvs {
		if ok, err := update(kv); ok {
			break
		} else if err != nil {
			return err
		}
	}
	go func() {
		watchCh := cli.Watch(context.Background(), key)
		for resp := range watchCh {
			for _, evt := range resp.Events {
				if evt.Type == etcdv3.EventTypePut {
					if ok, err := update(evt.Kv); ok {
						break
					} else if err != nil {
						break
					}
				}
			}
		}
	}()
	return nil
}

func GetClusterConfig() Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return *config
}

type Node struct {
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Proto   string `json:"proto"`
}

type cluster struct {
	sync.RWMutex
	cli     *etcdv3.Client
	service string
	once    *sync.Once
	deregCh map[string]chan struct{}
	nodes   map[string]*Node
}

var defaultCluster *cluster
var once = &sync.Once{}

func Init(service string) {
	once.Do(func() {
		defaultCluster = &cluster{
			cli:     etcd.GetClient(),
			service: service,
			once:    &sync.Once{},
			deregCh: make(map[string]chan struct{}),
			nodes:   make(map[string]*Node),
		}
	})
}

func Register(node *Node, ttl int) error {
	const minTTL = 2
	c := defaultCluster
	key := c.makeKey(node)
	if ttl < minTTL {
		ttl = minTTL
	}
	var errCh = make(chan error)
	go func() {
		kv := etcdv3.NewKV(c.cli)
		closeCh := make(chan struct{})
		lease := etcdv3.NewLease(c.cli)
		val, _ := json.Marshal(node)
		var curLeaseId etcdv3.LeaseID = 0
		ticker := time.NewTicker(time.Duration(ttl/2) * time.Second)
		register := func() error {
			if curLeaseId == 0 {
				leaseResp, err := lease.Grant(context.TODO(), int64(ttl))
				if err != nil {
					return err
				}
				if _, err := kv.Put(context.TODO(), key, string(val), etcdv3.WithLease(leaseResp.ID)); err != nil {
					return err
				}
				curLeaseId = leaseResp.ID
			} else {
				// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑
				if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
					curLeaseId = 0
				}
			}
			return nil
		}
		if err := register(); err != nil {
			logrus.Error("register node failed, error:", err)
			errCh <- err
		}
		close(errCh)
		for {
			select {
			case <-ticker.C:
				if err := register(); err != nil {
					logrus.Error("register node failed, error:", err)
					panic(err)
				}
			case <-closeCh:
				ticker.Stop()
				return
			}
		}
	}()
	err := <-errCh
	return err
}

func Deregister(node *Node) error {
	c := defaultCluster
	c.Lock()
	defer c.Unlock()
	key := c.makeKey(node)
	if ch, ok := c.deregCh[key]; ok {
		close(ch)
		delete(c.deregCh, key)
	}
	_, err := c.cli.Delete(context.Background(), key, etcdv3.WithPrefix())
	return err
}

func Discover() (output []*Node, err error) {
	c := defaultCluster
	key := fmt.Sprintf("/%s/nodes/", c.service)
	c.once.Do(func() {
		var resp *etcdv3.GetResponse
		resp, err = c.cli.Get(context.Background(), key, etcdv3.WithPrefix())
		if err != nil {
			return
		}
		for _, kv := range resp.Kvs {
			k := string(kv.Key)
			if len(k) > len(key) {
				var node *Node
				json.Unmarshal(kv.Value, &node)
				if node != nil {
					c.Lock()
					c.nodes[k] = node
					c.Unlock()
				}
			}
		}
		watchCh := c.cli.Watch(context.Background(), key, etcdv3.WithPrefix())
		go func() {
			for {
				select {
				case resp := <-watchCh:
					for _, evt := range resp.Events {
						k := string(evt.Kv.Key)
						if len(k) <= len(key) {
							continue
						}
						switch evt.Type {
						case etcdv3.EventTypePut:
							var node *Node
							json.Unmarshal(evt.Kv.Value, &node)
							if node != nil {
								c.Lock()
								c.nodes[k] = node
								c.Unlock()
							}
						case etcdv3.EventTypeDelete:
							c.Lock()
							if _, ok := c.nodes[k]; ok {
								delete(c.nodes, k)
							}
							c.Unlock()
						}
					}
				}
			}
		}()

	})
	if err != nil {
		return nil, err
	}
	c.RLock()
	for _, node := range c.nodes {
		output = append(output, node)
	}
	c.RUnlock()
	return
}

func (c *cluster) makeKey(node *Node) string {
	id := strings.Replace(node.Addr, ".", "-", -1)
	id = strings.Replace(id, ":", "-", -1)
	return fmt.Sprintf("/%s/nodes/%s", c.service, id)
}
