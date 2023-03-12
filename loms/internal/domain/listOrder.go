package domain

func (d *Domain) GetListOrder(orderID int64) (Order, error) {
	// Fixture
	return Order{
		New,
		123,
		[]Item{
			{
				1076963,
				2,
			},
			{
				2956315,
				5,
			},
			{
				1625903,
				1,
			},
		},
	}, nil
}
