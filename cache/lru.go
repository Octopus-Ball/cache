// @Description lru缓存的实现
// @Author zygoodest

package cache

import (
	"fmt"
	"sync"
	"time"
)

// *************************************************************************************************

func initLRUCache(typ string, cap int, duration time.Duration) *LRUCache {
	lruCache := &LRUCache{
		cap:      cap,
		duration: duration,
		rwm:      new(sync.RWMutex),
		li:       newList(cap),
		dic:      make(map[string]*listNode, cap),
	}
	lruCache.loop(duration / 2)
	return lruCache
}

// LRUCache 符合LUR算法的缓冲池
type LRUCache struct {
	cap      int
	duration time.Duration
	rwm      *sync.RWMutex
	li       *list
	dic      map[string]*listNode
}

// *************************************************************************************************

// Len 返回缓存内已有的数据量
func (c *LRUCache) Len() int {
	return c.li.len
}

// Cap 返回缓存的容量
func (c *LRUCache) Cap() int {
	return c.cap
}

// Duration 获取缓存默认有效时间
func (c *LRUCache) Duration() time.Duration {
	return c.duration
}

// *************************************************************************************************

// Get 根据key从缓存池获取值
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.rwm.RLock() // 读锁
	defer c.rwm.RUnlock()

	node, ok := c.dic[key]
	if !ok {
		fmt.Println("未命中缓存")
		return nil, false
	}

	// 若缓存过期了(删除该节点，并返回false)
	if isTimeout(node.deadline) {
		fmt.Println("缓存过期了，对其惰性删除")
		c.del(node)
		return nil, false
	}

	fmt.Println("缓存命中")
	// 若存在的值被使用到，则将其提升到链表头部
	c.li.upToHead(node)
	return node.data.data, true
}

// Put 将数据放入缓存池
func (c *LRUCache) Put(key string, data interface{}) {
	c.rwm.Lock() // 写锁
	defer c.rwm.Unlock()

	node := c.li.addNode(key, data)
	c.dic[key] = node

	c.limitLen()
}

// Del 从缓冲池删除数据（通过key）
func (c *LRUCache) Del(key string) {
	node, ok := c.dic[key]
	if !ok {
		return
	}
	c.del(node)
}

// *************************************************************************************************

// RegistGetFunc 将没有cache时获取数据的方法注册到缓存
func (c *LRUCache) RegistGetFunc(getNoCache func(string) interface{}) (getWithCache func(string) interface{}) {
	getWithCache = func(key string) interface{} {
		data, exist := c.Get(key)
		if !exist {
			data = getNoCache(key)
			c.Put(key, data)
		}

		return data
	}
	return
}

// RegistPutFunc 将没有cache时设置数据的方法注册到缓存
func (c *LRUCache) RegistPutFunc(putNoCache func(string, interface{})) (putWithCache func(string, interface{})) {
	putWithCache = func(key string, value interface{}) {
		// 先更新数据，再淘汰缓存
		putNoCache(key, value)
		c.Del(key)
	}
	return
}

// RegistDelFunc 将没有cache时删除数据的方法注册到缓存
func (c *LRUCache) RegistDelFunc(delNoCache func(string)) (delWithCache func(string)) {
	delWithCache = func(key string) {
		// 先删除数据，再淘汰缓存
		delNoCache(key)
		c.Del(key)
	}
	return
}

// *************************************************************************************************

// 从缓存池删除数据（通过node）
func (c *LRUCache) del(node *listNode) {
	// c.rwm.Lock() // 写锁
	// defer c.rwm.Unlock()
	// TODO ************************************* 后续考虑此处操作是否需要放到list里并加锁
	c.li.takeNode(node)
	delete(c.dic, node.data.key)
}

// 获得过期时间点
// 在标准过期时长基础上增减一定值以防止缓存雪崩)
func (c *LRUCache) getEnding() (deadline int64) {
	random := time.Now().UnixNano() % 100                                    // 0~99的随机数
	duration := c.duration + (c.duration / 500 * time.Duration(int(random))) // 添加随机数后的持续时间
	deadline = time.Now().Add(duration).UnixNano()                           // 最终的截至时间

	return
}

// 批量删除过期节点，从链表尾往前遍历，删除过期的节点(返回清除掉节点的数量)
func (c *LRUCache) cleanTimeoutNode() (cleanSum int) {
	for node := c.li.end.pre; node != c.li.head; node = node.pre {
		if isTimeout(node.deadline) {
			c.del(node)
			cleanSum++
		} else { // 直到遇到一个未过期的，则停止清理
			break
		}
	}
	return
}

// 若超过容量，则删除最后一个节点
func (c *LRUCache) limitLen() {
	if c.Len() < c.Cap() {
		fmt.Println("未超容量")
		return
	}

	lastNode := c.li.end.pre
	fmt.Println("超过了容量，则删除最后一个节点", lastNode.data.data)
	c.del(lastNode)
}

// 运行批量清除
func (c *LRUCache) runClean() {
	// 当现有节点数量大于容量的一半时，运行批量删除过期节点
	fmt.Println(c.Len())
	fmt.Println(c.Cap() / 2)
	if c.Len() < (c.Cap() / 2) {
		fmt.Println("无需进行批量清除")
		return
	}
	fmt.Println("需要运行批量清除")
	sum := c.cleanTimeoutNode()
	fmt.Println("运行批量清除完毕，删除掉节点", sum)
}

// 运行定时任务
func (c *LRUCache) loop(duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			c.runClean()
		}
	}()
}

// *************************************************************************************************

// 根据传入的到期时间戳判断是否到期
func isTimeout(deadline int64) bool {
	if time.Now().UnixNano() < deadline {
		return false
	}
	return true
}

// 已添加防止缓存雪崩的功能
// 待添加防止缓存击穿和缓存穿透的功能
