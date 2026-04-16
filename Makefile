MODULE=github.com/solardome/gamepulse-platform
PROTO_SRC=account/v1/account.proto

.PHONY: tidy generate generate-proto generate-graphql run-account run-accounts-gateway run-accounts-web compose-up compose-down

tidy:
	go mod tidy

generate: generate-proto generate-graphql

generate-proto:
	protoc \
		--proto_path=api/proto \
		--go_out=. \
		--go_opt=module=$(MODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(MODULE) \
		$(PROTO_SRC)

generate-graphql:
	go tool gqlgen generate --config accounts-gateway/gqlgen.yml

run-account:
	go run ./accounts-service/cmd/server

run-accounts-gateway:
	go run ./accounts-gateway/cmd/server

run-accounts-web:
	go run ./accounts-web/cmd/server

compose-up:
	docker compose up --build

compose-down:
	docker compose down --remove-orphans
