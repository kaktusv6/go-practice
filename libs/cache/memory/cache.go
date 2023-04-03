package memory

import (
	"route256/libs/cache"
	cacheMetrics "route256/libs/cache/metrics"
	"sync"
	"time"
)

type MemoryCache struct {
	sync.RWMutex
	mapKeyToItem     map[string]Item
	defaultExpiredAt int64
	cleanerInterval  time.Duration
}

func NewMemoryCache(defaultExpiredAt int64, cleanerInterval time.Duration) *MemoryCache {
	memoryCache := &MemoryCache{
		mapKeyToItem:     make(map[string]Item),
		defaultExpiredAt: defaultExpiredAt,
		cleanerInterval:  cleanerInterval,
	}
	memoryCache.runCleaner()
	return memoryCache
}

// runCleaner Метод запуска чистильщика. Если у кэша стоит cleanerInterval = 0 то кэш не будет подчищаться
func (m *MemoryCache) runCleaner() {
	if m.cleanerInterval != 0 {
		go m.cleaner()
	}
}

// cleaner Метод очистки кэша у которого исте TTL
func (m *MemoryCache) cleaner() {
	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(m.cleanerInterval)

		// Если
		if m.mapKeyToItem == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		m.DeleteMany(
			m.expiredKeys(),
		)
	}
}

func (m *MemoryCache) expiredKeys() []string {
	result := make([]string, 0, len(m.mapKeyToItem))
	for key, item := range m.mapKeyToItem {
		diffTimeStorage := time.Now().Unix() - item.CreatedAd.Unix()
		if diffTimeStorage >= item.Expiration {
			result = append(result, key)
		}
	}
	return result
}

// Get Метод получения кэша по ключу
func (m *MemoryCache) Get(key string) (interface{}, bool) {
	cacheMetrics.CacheRequestsTotal.Inc()

	// Делаем блокировку к данным
	m.RLock()
	defer m.RUnlock()

	// Получаем жанные из кэша
	item, ok := m.mapKeyToItem[key]

	// Если нет данных в кэше то возврааем nil
	if !ok {
		cacheMetrics.CacheMissCount.Inc()
		return nil, false
	} else {
		cacheMetrics.CacheHitCount.Inc()
	}

	// Если время хранения истекло у данных то также ничего не возвращаем
	if item.Expiration > 0 && time.Now().Unix()-item.CreatedAd.Unix() >= item.Expiration {
		return nil, false
	}

	// Возвращаем данные
	return item.Value, true
}

// Set Метод создания или изменения кэша по ключу
func (m *MemoryCache) Set(key string, value interface{}, options ...cache.Option) {
	var ttl int64

	for _, option := range options {
		switch option.Key {
		case cache.TTL:
			ttl = option.Value.(int64)
		}
	}

	if ttl == 0 {
		ttl = m.defaultExpiredAt
	}

	m.mapKeyToItem[key] = Item{
		Value:      value,
		Expiration: ttl,
		CreatedAd:  time.Now(),
	}
}

// Delete Метод удаления кэша из памяти
func (m *MemoryCache) Delete(key string) {
	delete(m.mapKeyToItem, key)
}

// GetMany Метод получения значений по набору ключей
func (m *MemoryCache) GetMany(keys []string) ([]interface{}, []bool) {
	keysLen := len(keys)
	results := make([]interface{}, keysLen, keysLen)
	resultChecks := make([]bool, keysLen, keysLen)

	for index, key := range keys {
		results[index], resultChecks[index] = m.Get(key)
	}

	return results, resultChecks
}

// SetMany Метод изменения кэша по набору ключей i-ый ключ соответствует i-му значению
func (m *MemoryCache) SetMany(keys []string, values []interface{}, options ...cache.Option) {
	for index, value := range values {
		m.Set(keys[index], value, options...)
	}
}

// DeleteMany Удаление кэша по нескольким ключам
func (m *MemoryCache) DeleteMany(keys []string) {
	for _, key := range keys {
		m.Delete(key)
	}
}
