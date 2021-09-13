package cache

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string, result interface{}) (interface{}, error)
	Del(key string) error
	FlushDB() error
}
