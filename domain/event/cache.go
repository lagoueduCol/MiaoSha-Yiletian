package event

import (
	"sync"

	"github.com/letian0805/seckill/infrastructure/stores"
)

type Cache struct {
	cache stores.ObjCache
}

var cache *Cache
var once = &sync.Once{}

func InitCache() error {
	var err error
	once.Do(func() {
		cache = &Cache{
			cache: stores.NewObjCache(),
		}
		err = cache.load()
	})
	return err
}

func GetCache() *Cache {
	return cache
}

func (c *Cache) load() error {
	return nil
}

func (c *Cache) GetList() *Topic {

	return nil
}

func (c *Cache) getCurrentEvent() *Event {

	return nil
}

func (c *Cache) GetEventInfo(goodsID int64) *Info {

	return nil
}
