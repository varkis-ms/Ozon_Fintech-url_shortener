ifeq ($(shell test -e '.env' && echo -n yes),yes)
	include .env
endif

HELP_FUN = \
	%help; while(<>){push@{$$help{$$2//'options'}},[$$1,$$3] \
	if/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/}; \
    print"$$_:\n", map"  $$_->[0]".(" "x(20-length($$_->[0])))."$$_->[1]\n",\
    @{$$help{$$_}},"\n" for keys %help; \
database = in-memory
# Commands
env:
	@$(eval SHELL:=/bin/bash)
	@cp .env.sample .env

help:
	@echo -e "Usage: make [target] ...\n"
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

build:
	go build -o ./cmd/url_shortener

run:
	go run ./cmd/url_shortener -database=$(database)

df:
	docker build --build-arg DATABASE=$(database) --tag shortener_url_service .

service_up:
	docker-compose -f docker-compose.service.yml up -d --remove-orphans

docker_db:
	docker-compose -f docker-compose.db.yml up -d --remove-orphans

docker_test:
	docker-compose -f docker-compose.test.yml up -d --remove-orphans

integration-test:
	go test -tags integration ./internal/tests/integration/...

unit-test:
	go test ./...

proto:
	rm -f pkg/pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pkg/pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger \
	--openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true \
	proto/shortener.proto

.PHONY: db proto df run build service_up docker_db docker_test unit-test integration-test