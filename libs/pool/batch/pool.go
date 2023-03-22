package batch

import (
	"context"
	"sync"
)

// Pool Интерфейс пула задач
type Pool[In, Out any] interface {
	Submit(context.Context, []Task[In, Out])
	Close()
	GetResultsChannel() <-chan Out
}

type pool[In, Out any] struct {
	// Кол-во воркеров которое будет запущено
	amountWorkers int

	// Счетчик активных воркеров
	workersWaitGroup sync.WaitGroup

	// Канал полчения задач
	tasksSource chan Task[In, Out]
	// Канал с итоговыми результатами выполненых задач
	tasksResultSource chan Out
}

func NewPool[In, Out any](
	ctx context.Context,
	amountWorkers int,
) Pool[In, Out] {
	pool := &pool[In, Out]{
		amountWorkers: amountWorkers,
	}

	// Настраиваем Pool
	pool.bootstrap(ctx)

	return pool
}

func (p *pool[In, Out]) GetResultsChannel() <-chan Out {
	return p.tasksResultSource
}

// Close Закрываем все каналы
func (p *pool[In, Out]) Close() {
	// Больше задач не будет
	close(p.tasksSource)

	// Ожидаем окончания раболты всех воркеров
	p.workersWaitGroup.Wait()

	// Закрываем канал выдачи результата задач
	close(p.tasksResultSource)
}

// Submit Отправляем все задачи в канал на выполнение воркерам
func (p *pool[In, Out]) Submit(ctx context.Context, tasks []Task[In, Out]) {
	go func() {
		for _, task := range tasks {
			p.tasksSource <- task
		}
		p.Close()
	}()
}

func (p *pool[In, Out]) bootstrap(ctx context.Context) {
	// Формируем каналы пула
	p.tasksSource = make(chan Task[In, Out], p.amountWorkers)
	p.tasksResultSource = make(chan Out, p.amountWorkers)

	// Запускаем указанное кол-во воркеров
	for i := 0; i < p.amountWorkers; i++ {
		p.workersWaitGroup.Add(1)
		worker := NewWorker[In, Out](
			i+1,
			p.tasksSource,
			p.tasksResultSource,
			&p.workersWaitGroup,
		)
		worker.Run(ctx)
	}
}
