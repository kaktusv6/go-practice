package batch

import (
	"context"
	"sync"
)

type Worker interface {
	Run(context.Context)
}

type worker[In, Out any] struct {
	index        int
	tasksSource  chan Task[In, Out]
	resultSource chan Out
	wg           *sync.WaitGroup
}

func NewWorker[In, Out any](
	index int,
	taskSource chan Task[In, Out],
	resultSource chan Out,
	wg *sync.WaitGroup,
) Worker {
	return &worker[In, Out]{
		index:        index,
		tasksSource:  taskSource,
		resultSource: resultSource,
		wg:           wg,
	}
}

func (w *worker[In, Out]) Run(ctx context.Context) {
	go func() {
		defer w.wg.Done()

		for task := range w.tasksSource {
			select {
			case <-ctx.Done():
				return
			case w.resultSource <- task.Callback(task.InputArgs):
			}
		}
	}()
}
