package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	mu             sync.Mutex                    // A mutex to ensure thread-safe operations
	capacity       int                           // Cache capacity
	storage        map[interface{}]*list.Element // hashmap: key-item key, value doubly linked list node reference
	evictionList   *list.List                    // doubly linked list to keep track of the order of items for eviction purposes
	evictionPolicy EvictionPolicy                // An interface to define the eviction policy (e.g., FIFO, LRU, LIFO)
}

// CacheItem is represents a cached item
// for internal use only
type CacheItem struct {
	Key   interface{} // The key of the cache item
	Value interface{} // The value of the cache item
}

type EvictionPolicy interface {
	Add(evictionList *list.List, item *list.Element)    // Method to add an item to the eviction list
	Remove(evictionList *list.List) *list.Element       // Method to remove an item from the eviction list
	Access(evictionList *list.List, item *list.Element) // Method to mark an item as accessed
}

// A constructor function to create a new Cache instance
// Initializes the capacity, storage, evictionList, and evictionPolicy fields
// Returns a pointer to the new Cache instance
func NewCache(capacity int, policy EvictionPolicy) *Cache {
	return &Cache{
		capacity:       capacity,
		storage:        make(map[interface{}]*list.Element),
		evictionList:   list.New(),
		evictionPolicy: policy,
	}
}

func (c *Cache) Set(key, value interface{}) {
	//  mutex lock
	// Prevents multiple threads from writing into to the same key in parallel
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the item is already in the cache
	// If item is already in cache it removes  it from the cache
	if el, ok := c.storage[key]; ok {
		c.evictionList.Remove(el)
		delete(c.storage, key)
	}

	// Check cache has the capacity to store the new item
	// If the capacity is full delete the item based on the eviction policy
	if len(c.storage) >= c.capacity {
		el := c.evictionPolicy.Remove(c.evictionList)
		if el != nil {
			item := el.Value.(*CacheItem)
			delete(c.storage, item.Key)
			c.evictionList.Remove(el)
		}
	}

	// If the item is not in cache add it to the cache (hashmap and doubly linked list) after the capacity has been checked
	item := &CacheItem{Key: key, Value: value}
	el := c.evictionList.PushFront(item)
	c.storage[key] = el
	c.evictionPolicy.Add(c.evictionList, el)
}

// Get returns the cached element corresponding to the given key
// It also calls the access function of the given eviction policy(only required for LRU to move the item to the front of the list)
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.storage[key]; ok {
		c.evictionPolicy.Access(c.evictionList, el)
		return el.Value.(*CacheItem).Value, true
	}
	return nil, false
}

// FIFO (First In First Out)
type FIFO struct{}

func NewFIFO() *FIFO {
	return &FIFO{}
}

func (p *FIFO) Add(evictionList *list.List, item *list.Element) {
	// No operation needed for FIFO add
}

func (p *FIFO) Remove(evictionList *list.List) *list.Element {
	// FIFO removes from the back (oldest item)
	return evictionList.Back()
}

func (p *FIFO) Access(evictionList *list.List, item *list.Element) {
	// No operation needed for FIFO access
}

// LRU (Least Recently Used)
type LRU struct{}

func NewLRU() *LRU {
	return &LRU{}
}

func (p *LRU) Add(evictionList *list.List, item *list.Element) {
	// No operation needed for LRU add
	// Note: Least recently used item will be at the back of the doubly linked list(last node in doubly linked list)

}

func (p *LRU) Remove(evictionList *list.List) *list.Element {
	// LRU removes from the back (least recently used item)
	return evictionList.Back()
}

func (p *LRU) Access(evictionList *list.List, item *list.Element) {
	// Moves the item to the front of the list
	evictionList.MoveToFront(item)
}

// LIFO (Last In First Out)
type LIFO struct{}

func NewLIFO() *LIFO {
	return &LIFO{}
}

func (p *LIFO) Add(evictionList *list.List, item *list.Element) {
	// No operation needed for LIFO add
}

func (p *LIFO) Remove(evictionList *list.List) *list.Element {
	// LIFO removes from the front (most recently added item)
	return evictionList.Front()
}

func (p *LIFO) Access(evictionList *list.List, item *list.Element) {
	// No operation needed for LIFO access
}
