include .env
export

install:
	# install dependencies
	go mod tidy

docker-up:
	docker-compose -f ./docker/docker-compose-local.yml up -d

docker-down:
	docker-compose -f ./docker/docker-compose-local.yml down -v

amqp:
	go run ./cmd start amqp

.PHONY: install docker-up docker-down amqp
