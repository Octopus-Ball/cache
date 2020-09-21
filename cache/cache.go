// @Description 各种缓存的实现
// @Author zygoodest

package cache

import (
	"fmt"
	"time"
)

const (
	// TypeTimeOut 时效缓存
	TypeTimeOut = "timeout"
	// TypeLRU LRU缓存
	TypeLRU = "lru"
	// TypeLFU LFU缓存
	TypeLFU = "lfu"
)

// InitCache 实例化一个缓存对象
// @param     	typ        string         	缓存类型
// @param     	cap        int         		缓存容量
// @param     	timeOut    time.Duration 	缓存过期时间
// @return    	c			*Cache			缓存实例
func InitCache(typ string, cap int, duration time.Duration) (c Cache) {
	switch typ {
	case TypeLRU:
		c = initLRUCache(typ, cap, duration)
	default:
		fmt.Println("类型不存在")
	}

	return
}

// Cache 缓存池
type Cache interface {
	// @description 获取当前被缓存数据的量
	// @return		len			int				已缓存数据的量
	Len() (len int)
	// @description 获取缓存的最大容量
	// @return		cap			int				缓存的容量
	Cap() (cap int)
	// @description 获取默认的缓存过期时间
	// @return		d			time.Duration	默认缓存过期时间
	Duration() (d time.Duration)
	// @description 通过key从缓存获取value
	// @param     	key			string         	key值
	// @return		v			interface{}		value值
	Get(key string) (v interface{}, exist bool)
	// @description 通过key-value的形式往缓存添加/更新数据
	// @param     	key			string         	key值
	// @param     	value		interface{}     value值
	Put(key string, value interface{})
	// @description 通过key删掉缓存里对应的数据
	Del(key string)
	// @description 将无cache获取数据的方式注册到缓存，以获得带缓存的获取数据的方式
	RegistGetFunc(getNoCache func(string) interface{}) (getWithCache func(string) interface{})
	// @description 将无cache添加数据的方式注册到缓存，以获得带缓存的添加数据的方式
	RegistPutFunc(putNoCache func(string, interface{})) (putWithCache func(string, interface{}))
	// @description 将无cache删除数据的方式注册到缓存，以获得带缓存的删除数据的方式
	RegistDelFunc(delNoCache func(string)) (delWithCache func(string))
}

// *************************************************************************************************
