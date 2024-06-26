# Cache Library

## Overview

This is an in-memory caching library written in Go, designed to support multiple standard eviction policies including FIFO (First In First Out), LIFO (Last In First Out), and LRU (Least Recently Used). The library also allows the addition of custom eviction policies and is thread-safe.

## Features

- **Multiple Eviction Policies**: Supports FIFO, LIFO, and LRU.
- **Custom Eviction Policies**: Allows adding custom eviction policies.
- **Thread Safety**: Ensures thread-safe operations using mutexes.
- **Extensible Design**: The low-level design is extensible for future enhancements.

## Installation

To install the package, use the following command:

```bash
go get github.com/atul-007/cached
```

## Usage

Below is a sample usage of the caching library demonstrating all the eviction policies.

### Import the Package

```go
import (
    "fmt"
    "github.com/atul-007/cached"
)
```

### Example Main Function

```go
func testCache(policyName string, c *cache.Cache) {
	fmt.Printf("Testing %s Cache\n", policyName)

	// Initial adds
	fmt.Println("Adding: (1, 'one')")
	c.Set(1, "one")
	fmt.Println("Adding: (2, 'two')")
	c.Set(2, "two")
	fmt.Println("Adding: (3, 'three')")
	c.Set(3, "three")
	fmt.Println("Cache state after initial adds (should contain keys 1, 2, 3):")
	for i := 1; i <= 3; i++ {
		value, ok := c.Get(i)
		if ok {
			fmt.Printf("Get: %d -> %s\n", i, value)
		} else {
			fmt.Printf("Get: %d -> not found\n", i)
		}
	}

	// Access some elements
	fmt.Println("Accessing key 1")
	c.Get(1) // Access key 1

	// Add more elements to trigger eviction
	fmt.Println("Adding: (4, 'four')")
	c.Set(4, "four")
	fmt.Println("Cache state after adding key 4 (should have evicted one item):")
	for i := 1; i <= 4; i++ {
		value, ok := c.Get(i)
		if ok {
			fmt.Printf("Get: %d -> %s\n", i, value)
		} else {
			fmt.Printf("Get: %d -> not found\n", i)
		}
	}

	// Add another element to further test eviction
	fmt.Println("Adding: (5, 'five')")
	c.Set(5, "five")
	fmt.Println("Cache state after adding key 5 (should have evicted another item):")
	for i := 1; i <= 5; i++ {
		value, ok := c.Get(i)
		if ok {
			fmt.Printf("Get: %d -> %s\n", i, value)
		} else {
			fmt.Printf("Get: %d -> not found\n", i)
		}
	}

	fmt.Println()
}

func main() {
	fifoCache := cache.NewCache(3, cache.NewFIFO())
	testCache("FIFO", fifoCache)

	lruCache := cache.NewCache(3, cache.NewLRU())
	testCache("LRU", lruCache)

	lifoCache := cache.NewCache(3, cache.NewLIFO())
	testCache("LIFO", lifoCache)
}
```

## Detailed Functionality

### Set Function

The `Set` function adds a key-value pair to the cache. If the key already exists, it updates the value and refreshes its position based on the eviction policy. If the cache exceeds its capacity, it evicts an element based on the eviction policy.

### Get Function

The `Get` function retrieves the value associated with a key. If the key exists, it updates the key's position based on the eviction policy (for LRU).

## Eviction Policies

### FIFO (First In First Out)

FIFO evicts the oldest element (the one that was added first).

### LIFO (Last In First Out)

LIFO evicts the newest element (the one that was added last).

### LRU (Least Recently Used)

LRU evicts the least recently used element. When an element is accessed, it is moved to the front to mark it as recently used.

## Custom Eviction Policies

You can add custom eviction policies by implementing the `EvictionPolicy` interface:

```go
type EvictionPolicy interface {
    Add(item *list.Element)
    Remove() *list.Element
    Access(item *list.Element)
}
```
## Custom Eviction Policy Example

A custom eviction policy called Random Eviction Policy where a random item is removed from the cache when it reaches its capacity.

## Random Eviction Policy Implementation

```go
package cache

import (
	"container/list"
	"math/rand"
)

// RandomEvictionPolicy struct to implement the custom eviction policy
type RandomEvictionPolicy struct{}

// NewRandomEvictionPolicy constructor for RandomEvictionPolicy
func NewRandomEvictionPolicy() *RandomEvictionPolicy {
	return &RandomEvictionPolicy{}
}

// Add method for RandomEvictionPolicy
func (p *RandomEvictionPolicy) Add(evictionList *list.List, item *list.Element) {
	// No operation needed for random eviction add
}

// Remove method for RandomEvictionPolicy
func (p *RandomEvictionPolicy) Remove(evictionList *list.List) *list.Element {
	// Randomly select an element to remove
	n := rand.Intn(evictionList.Len())
	for e := evictionList.Front(); e != nil; e = e.Next() {
		if n == 0 {
			return e
		}
		n--
	}
	return nil
}

// Access method for RandomEvictionPolicy
func (p *RandomEvictionPolicy) Access(evictionList *list.List, item *list.Element) {
	// No operation needed for random eviction access
}

```
## Contributing

Feel free to fork the repository and submit pull requests. For major changes, please open an issue first to discuss what you would like to change.

---

