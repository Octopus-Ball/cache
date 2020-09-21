package main

import (
	"time"
	"cache/cache"
	"fmt"
)

func getNoCache(key string) (value interface{}) {
	fmt.Println("调用外存")
	value = key
	return
}

func main() {
	c := cache.InitCache(cache.TypeLRU, 5, time.Second)
	c.Put("a", "aa")
	c.Put("b", "bb")
	c.Put("c", "cc")
	c.Put("d", "dd")
	c.Put("e", "ee")

	time.Sleep(time.Second * 100)


	// getWithCache := c.RegistGetFunc(getNoCache)
	// getWithCache("aa")
	// getWithCache("bb")
	// v := getWithCache("bb")
	// fmt.Println(v)
}
