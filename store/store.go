package store

//KVStore is the interface for basic implementation of the kv store
type KVStore interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string)
	PrintEntries() string
	GetStats() Stats
}

//Stats encapsulates the statistics
type Stats struct {
	Sets       int
	Gets       int
	Success    int
	MemoryUsed int
	Capacity   int
}
