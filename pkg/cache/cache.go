package cache

var InMemory *InMemoryDBImpl

type (
	InMemoryDB interface {
		Get(key string) (interface{}, bool)
		Set(key string, value interface{})
		Delete(key string)
		Clear()
	}
	InMemoryDBImpl struct {
		cache map[string]interface{}
	}
)

// NewInMemoryDB creates a new in-memory cache.
func NewInMemoryDB() *InMemoryDBImpl {
	return &InMemoryDBImpl{
		cache: make(map[string]interface{}),
	}
}

// Get returns the value for the given key.
func (db *InMemoryDBImpl) Get(key string) (interface{}, bool) {
	value, ok := db.cache[key]
	return value, ok
}

// Set sets the value for the given key.
func (db *InMemoryDBImpl) Set(key string, value interface{}) {
	db.cache[key] = value
}

// Delete deletes the value for the given key.
func (db *InMemoryDBImpl) Delete(key string) {
	delete(db.cache, key)
}

// Clear clears the cache.
func (db *InMemoryDBImpl) Clear() {
	db.cache = make(map[string]interface{})
}
