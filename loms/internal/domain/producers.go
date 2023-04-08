package domain

//go:generate bash -c "rm -rf ./mocks/order_status_notifier_minimock.go"
//go:generate bash -c "mkdir -p mocks"
//go:generate minimock -i OrderStatusNotifier -o ./mocks/ -s "_minimock.go"

type OrderStatusNotifier interface {
	Notify(order *Order) error
}
