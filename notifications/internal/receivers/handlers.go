package receivers

func InitHandlers() map[string]ConsumerHandlerFunc {
	return map[string]ConsumerHandlerFunc{
		"order_statuses": orderStatusHandler,
	}
}
