package repositories

import (
	"context"
	"sync"
	"time"
)

import (
	"golang.org/x/time/rate"
	"route256/checkout/internal/domain"
	productServiceV1Clinet "route256/checkout/pkg/product_service_v1"
	"route256/libs/pool/batch"
)

type TaskInputArgs struct {
	Index int
	Sku   uint32
}

type TaskOutResult struct {
	Index       int
	ProductInfo *domain.ProductInfo
}

type ProductRepository struct {
	productServiceClient productServiceV1Clinet.ProductServiceClient
	productServiceToken  string
}

func NewOrderProductRepository(
	productServiceClient productServiceV1Clinet.ProductServiceClient,
	productServiceToken string,
) domain.ProductRepository {
	return &ProductRepository{
		productServiceClient,
		productServiceToken,
	}
}

const (
	AmountGetterProductWorkers = 5
)

// GetListBySkus Метод получения списка продуктов на основе списка sku
// В методе используется pool worker и limiter rate для оптимального получения товаров от сервиса
func (p *ProductRepository) GetListBySkus(
	ctx context.Context,
	skus []uint32,
) ([]*domain.ProductInfo, error) {
	tasks := make([]batch.Task[*TaskInputArgs, *TaskOutResult], 0, len(skus))

	// Создаем лимитер на выполнение запросов к сервису
	limiter := rate.NewLimiter(rate.Every(time.Second), 20)

	for index, sku := range skus {
		// Формируем таски на получение информации о заказе
		tasks = append(tasks, batch.Task[*TaskInputArgs, *TaskOutResult]{
			Callback: func(args *TaskInputArgs) *TaskOutResult {
				limiter.Wait(ctx)

				productInfo, _ := p.GetProductBySku(ctx, args.Sku)

				return &TaskOutResult{
					Index:       args.Index,
					ProductInfo: productInfo,
				}
			},
			InputArgs: &TaskInputArgs{
				Index: index,
				Sku:   sku,
			},
		})
	}

	// Создаем pool worker
	pool := batch.NewPool[*TaskInputArgs, *TaskOutResult](ctx, AmountGetterProductWorkers)

	// Создаем WaitGroup для сбора данных
	var wg sync.WaitGroup

	productInfoList := make([]*domain.ProductInfo, len(skus), len(skus))
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Собираем результат выполнения тасков
		for result := range pool.GetResultsChannel() {
			productInfoList[result.Index] = result.ProductInfo
		}
	}()

	// Отправляем таски в pool worker
	pool.Submit(ctx, tasks)

	// Ожидаем сбор всех данных
	wg.Wait()

	return productInfoList, nil
}

func (p *ProductRepository) GetProductBySku(ctx context.Context, sku uint32) (*domain.ProductInfo, error) {
	productResponse, err := p.productServiceClient.GetProduct(ctx, &productServiceV1Clinet.GetProductRequest{
		Token: p.productServiceToken,
		Sku:   sku,
	})

	if err != nil {
		return nil, err
	}

	return &domain.ProductInfo{
		Name:  productResponse.GetName(),
		Price: productResponse.GetPrice(),
	}, nil
}
