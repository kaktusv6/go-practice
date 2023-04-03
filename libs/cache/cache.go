package cache

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, options ...Option)
	Delete(key string)

	GetMany(keys []string) ([]interface{}, []bool)
	SetMany(keys []string, values []interface{}, options ...Option)
	DeleteMany(keys []string)
}

type Option struct {
	Key   OptionKey
	Value interface{}
}

type OptionKey string

const (
	TTL OptionKey = "ttl"
)

func OptionTTL(value int64) Option {
	return Option{Key: TTL, Value: value}
}
