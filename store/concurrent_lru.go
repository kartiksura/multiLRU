package store

import (
	"hash/crc32"
	"runtime"
)

//ConcurrentLRU is an LRU desingned for reducing lock contention
type ConcurrentLRU struct {
	stores []LRU
	nC     int //concurrency
}

//InitConcurrentLRU Initializes the LRU
func InitConcurrentLRU(capacity int) *ConcurrentLRU {
	var b ConcurrentLRU
	b.nC = runtime.NumCPU()

	for i := 0; i < b.nC; i++ {
		b.stores = append(b.stores, *InitLRU(capacity))
	}
	return &b
}

func (c *ConcurrentLRU) bucket(key string) int {
	checksum := int(crc32.ChecksumIEEE([]byte(key)))
	return (checksum % c.nC)
}

//Set  an item to the cache overwriting existing one if it
func (c *ConcurrentLRU) Set(key string, value []byte) error {
	return c.stores[c.bucket(key)].Set(key, value)
}

//Get returns the value if present
func (c *ConcurrentLRU) Get(key string) ([]byte, error) {
	return c.stores[c.bucket(key)].Get(key)
}

//Delete deletes the key
func (c *ConcurrentLRU) Delete(key string) {
	c.stores[c.bucket(key)].Delete(key)

}

//PrintEntries Prints the kv pairs
func (c *ConcurrentLRU) PrintEntries() string {
	ans := ""
	for i := 0; i < c.nC; i++ {
		ans = ans + c.stores[i].PrintEntries()
	}
	return ans
}

//GetStats gives the stats of the kv store
func (c *ConcurrentLRU) GetStats() Stats {
	var s Stats
	for i := 0; i < c.nC; i++ {
		t := c.stores[i].GetStats()
		s.Capacity += t.Capacity
		s.Gets += t.Gets
		s.MemoryUsed += t.MemoryUsed
		s.Sets += t.MemoryUsed
		s.Success += t.MemoryUsed
	}
	return s
}
