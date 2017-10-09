package store

import (
	"container/heap"
	"fmt"
	"log"
	"sync"
	"time"
)

type entry struct {
	key      string
	value    []byte    //
	accessed time.Time // time when the item is expired. it's okay to be stale.
	index    int       // index for priority queue needs. -1 if entry is free
	nSets    int
	nGets    int
	nHits    int
}

//LRU implments a lru cache
type LRU struct {
	lock        sync.Mutex
	table       map[string]*entry // all entries in table must be in lruList
	pq          priorityQueue     // some elements from table may be in priorityQueue
	capacity    int
	currentSize int
	stats       Stats
}

//InitLRU Initializes the LRU
func InitLRU(capacity int) *LRU {
	var b LRU
	b.table = make(map[string]*entry)
	b.capacity = capacity
	heap.Init(&b.pq)
	return &b
}

//Set  an item to the cache overwriting existing one if it
func (b *LRU) Set(key string, value []byte) error {
	log.Println("Setting:", key, string(value))
	b.lock.Lock()
	defer b.lock.Unlock()
	b.stats.Sets++
	if len(value) > b.capacity {
		return fmt.Errorf("The object is greater than the capacity")
	}
	if b.currentSize+len(value) > b.capacity {
		if err := b.evict(len(value)); err != nil {
			return err
		}
	}

	var e entry
	e.key = key
	e.value = value
	e.accessed = time.Now()
	if _, ok := b.table[key]; ok == true {
		b.currentSize -= len(e.value) //+ int(unsafe.Sizeof(e))
	}
	b.table[key] = &e
	b.currentSize += len(value) //+ int(unsafe.Sizeof(e))
	heap.Push(&b.pq, &e)
	log.Print("capacity:", b.capacity, " current size:", b.currentSize)
	return nil
}

func (b *LRU) evict(amt int) error {
	for amt > 0 {
		var e *entry
		e = heap.Pop(&b.pq).(*entry)
		log.Println("Evicted:", e.key, string(e.value))
		amt -= len(e.value) // + int(unsafe.Sizeof(e))
		b.currentSize -= len(e.value)
		delete(b.table, e.key)
		log.Print("after eviction capacity:", b.capacity, " current size:", b.currentSize)

	}
	return nil
}

//Get returns the value if present
func (b *LRU) Get(key string) ([]byte, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.stats.Gets++
	e, ok := b.table[key]
	if ok == true {
		b.table[key].accessed = time.Now()
		heap.Fix(&b.pq, e.index)
		b.stats.Success++
		return e.value, nil
	}
	return nil, fmt.Errorf("Key not found")

}

//Delete deletes the key
func (b *LRU) Delete(key string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	e, ok := b.table[key]
	if ok == true {
		b.currentSize -= len(e.value)
		delete(b.table, key)
		heap.Fix(&b.pq, e.index)

	}

}

//PrintEntries Prints the kv pairs
func (b *LRU) PrintEntries() string {
	b.lock.Lock()
	defer b.lock.Unlock()

	ans := ""
	for k, v := range b.table {
		ans = ans + "K: " + k + " V: " + string(v.value) + "\n"
	}
	return ans
}

//GetStats gives the stats of the kv store
func (b *LRU) GetStats() Stats {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.stats.MemoryUsed = b.currentSize
	b.stats.Capacity = b.capacity
	return b.stats
}
