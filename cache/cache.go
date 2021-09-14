package cache

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string, result interface{}) (interface{}, error)
	Del(keys ...string) error
	FlushDB() error
}
