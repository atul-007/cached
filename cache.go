package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	mu             sync.Mutex
	capacity       int
	storage        map[interface{}]*list.Element
	evictionList   *list.List
	evictionPolicy EvictionPolicy
}

type CacheItem struct {
	Key   interface{}
	Value interface{}
}

type EvictionPolicy interface {
	Add(item *list.Element)
	Remove() *list.Element
	Access(item *list.Element)
}

func NewCache(capacity int, policy EvictionPolicy) *Cache {
	return &Cache{
		capacity:       capacity,
		storage:        make(map[interface{}]*list.Element),
		evictionList:   list.New(),
		evictionPolicy: policy,
	}
}

func (c *Cache) Set(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.storage[key]; ok {
		c.evictionList.Remove(el)
		delete(c.storage, key)
	}

	if len(c.storage) >= c.capacity {
		el := c.evictionPolicy.Remove()
		item := el.Value.(*CacheItem)
		delete(c.storage, item.Key)
		c.evictionList.Remove(el)
	}

	item := &CacheItem{Key: key, Value: value}
	el := c.evictionList.PushFront(item)
	c.storage[key] = el
	c.evictionPolicy.Add(el)
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.storage[key]; ok {
		c.evictionPolicy.Access(el)
		return el.Value.(*CacheItem).Value, true
	}
	return nil, false
}

//First In First Out

type FIFO struct {
	list *list.List
}

func NewFIFO() *FIFO {
	return &FIFO{list: list.New()}
}

func (p *FIFO) Add(item *list.Element) {
	p.list.PushBack(item)
}

func (p *FIFO) Remove() *list.Element {
	return p.list.Remove(p.list.Front()).(*list.Element)
}

func (p *FIFO) Access(item *list.Element) {
	// No operation needed for FIFO
}
