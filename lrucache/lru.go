//go:build !solution

package lrucache

import (
	"container/list"
)

type LruCache struct {
	Cache           map[int]int
	KeyElementCache map[int]*list.Element
	LinkedList      *list.List
	Capacity        int
}

func New(cap int) LruCache {
	return LruCache{
		Capacity:        cap,
		Cache:           make(map[int]int, cap),
		KeyElementCache: make(map[int]*list.Element, cap),
		LinkedList:      list.New(),
	}
}

func (c *LruCache) Get(key int) (int, bool) {
	value, exist := c.Cache[key]
	if exist {
		c.touch(key)
		return value, true
	}
	return value, false
}

func (c *LruCache) touch(key int) {
	c.LinkedList.MoveToFront(c.KeyElementCache[key])
}

func (c *LruCache) touchByElement(e *list.Element) {
	c.LinkedList.MoveToFront(e)
}

func (c *LruCache) Set(key, value int) {
	if c.Capacity == 0 {
		return
	}

	if element, ok := c.KeyElementCache[key]; ok {
		c.touchByElement(element)
		c.Cache[key] = value
		return
	}

	if len(c.Cache) >= c.Capacity {
		oldestKey := c.LinkedList.Back().Value.(int)
		delete(c.KeyElementCache, oldestKey)
		delete(c.Cache, oldestKey)
		c.LinkedList.Remove(c.LinkedList.Back())
	}

	c.KeyElementCache[key] = c.LinkedList.PushFront(key)
	c.Cache[key] = value
}

func (c *LruCache) Range(f func(key, value int) bool) {
	for i := c.LinkedList.Back(); i != nil; i = i.Prev() {
		key := i.Value.(int)
		if !f(key, c.Cache[key]) {
			break
		}
	}
}

func (c *LruCache) Clear() {
	c.LinkedList = list.New()
	c.Cache = make(map[int]int, c.Capacity)
	c.KeyElementCache = make(map[int]*list.Element, c.Capacity)
}
