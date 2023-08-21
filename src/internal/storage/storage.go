package storage

// MemoryDB is a type alias for in memory-db type.
type MemoryDB map[string]any

// Storer defines storage behaviours.
type Storer interface {
	Set(key string, value any) any
	Get(key string) (any, error)
	Update(key string, value any) (any, error)
	Delete(key string) error
	List() MemoryDB
}
