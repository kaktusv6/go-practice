
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
