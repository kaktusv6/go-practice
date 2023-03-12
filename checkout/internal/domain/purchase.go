package domain

import "route256/libs/lomsClient"

func (d *Domain) Purchase(user int64) error {
	// Fixture
	fixtureItems := []lomsClient.Item{
		{1076963, 1},
		{1148162, 4},
		{1625903, 2},
	}

	_, err := d.lomsClient.CreateOrder(user, fixtureItems)
	return err
}
