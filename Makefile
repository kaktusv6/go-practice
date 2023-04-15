
build-all:
	cd checkout && GOOS=linux make build
	cd loms && GOOS=linux make build
	cd notifications && GOOS=linux make build

up-local-databases:
	docker-compose up -postgres-loms postgres-checkout

run-all: build-all
	docker-compose up --force-recreate --build checkout loms

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

high-load-checkout:
	hey call -d 1h localhost:50051 checkout_v1.CheckoutV1/ListCart 'user: 1'

high-load-loms:
	hey call -n 100 -d 1h localhost:50052 loms_v1.LomsV1/Stocks 'sku: 4996014'
