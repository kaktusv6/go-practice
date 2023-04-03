package rep_decorators

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/libs/cache"
	"strconv"
	"time"
)

type CacheProducts struct {
	repository domain.ProductRepository
	cache      cache.Cache
}

func NewCacheProducts(repository domain.ProductRepository, cache cache.Cache) domain.ProductRepository {
	return &CacheProducts{
		repository: repository,
		cache:      cache,
	}
}

func (c *CacheProducts) makeKeyBySku(sku uint32) string {
	return strconv.FormatUint(uint64(sku), 10)
}

func (c *CacheProducts) makeKeysBySkus(skus []uint32) []string {
	result := make([]string, 0, len(skus))
	for _, sku := range skus {
		result = append(result, c.makeKeyBySku(sku))
	}
	return result
}

func (c *CacheProducts) GetListBySkus(ctx context.Context, skus []uint32) ([]*domain.ProductInfo, error) {
	// Получаем данные из кэша
	products, cacheResults := c.cache.GetMany(c.makeKeysBySkus(skus))

	// То что получили от кша помещяем в результирующий массив
	result := make([]*domain.ProductInfo, 0, len(skus))
	for index, product := range products {
		if cacheResults[index] {
			result = append(result, product.(*domain.ProductInfo))
		}
	}

	// Формируем список sku продуктов который не пришел из кэша
	missingSkus := make([]uint32, 0, len(skus))
	for index, isGetOk := range cacheResults {
		if !isGetOk {
			missingSkus = append(missingSkus, skus[index])
		}
	}

	// Если список sku продуктов не пустой то выполняем запрос на получение у другого репозитория и сохраняем эти данные в кэше
	if len(missingSkus) > 0 {
		productsFromRepository, err := c.repository.GetListBySkus(ctx, missingSkus)
		if err != nil {
			return nil, err
		}

		productsForCache := make([]interface{}, 0, len(productsFromRepository))
		for _, product := range productsFromRepository {
			productsForCache = append(productsForCache, product)
			result = append(result, product)
		}

		// Сохраняем все данные на час
		c.cache.SetMany(
			c.makeKeysBySkus(missingSkus),
			productsForCache,
			cache.OptionTTL(int64(time.Hour)),
		)
	}

	// Взвращаем данные
	return result, nil
}

func (c *CacheProducts) GetProductBySku(ctx context.Context, sku uint32) (*domain.ProductInfo, error) {
	// ФОрмируем ключ для кэша
	key := c.makeKeyBySku(sku)

	var result *domain.ProductInfo

	// Получаем данные из кэша
	resultCache, ok := c.cache.Get(key)
	// Если данных нет, то запрашиваем их у репозитория, иначе возвращаем данные
	if !ok {
		result, err := c.repository.GetProductBySku(ctx, sku)
		if err != nil {
			return nil, err
		}

		// Сохраняем в кэше то что получили от репозитория
		c.cache.Set(key, result)
	} else {
		result = resultCache.(*domain.ProductInfo)
	}

	return result, nil
}
