package batch

// Task Задача попадающая в pool на выполненеи
type Task[In, Out any] struct {
	// Функция которая должна выполниться в рамках Task
	Callback func(In) Out
	// Аргументы функции
	InputArgs In
}
