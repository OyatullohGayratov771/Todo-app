gen-task:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./protos/task/task.proto


run:
	@go run cmd/main.go


CURRENT_DIR := $(shell pwd)
DB_URL := "postgres://postgres:0101@localhost:5432/todo_app?sslmode=disable"

migrate-up:
	migrate -path migrations/ -database $(DB_URL) up

migrate-down:
	migrate -path migrations/ -database  $(DB_URL) down

migrate-force:
	migrate -path migrations/ -database  $(DB_URL) force 1