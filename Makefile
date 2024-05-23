run:
	go run cmd/main.go

build:
	go build -o bin/bookstore cmd.main.go

export POSTGRES_DB=postgres://postgres:vakhaboff@localhost:5432/shaxboz?sslmode=disable

migrate-file:
	migrate create -ext sql -dir migrations/ -seq migrate_file

migrate-up:
	migrate -path pkg/migrations -database $(POSTGRES_DB) up

migrate-down:
	migrate -path migrations -database $(POSTGRES_DB) down

migrate-force:
	migrate -path pkg/migrations -database $(POSTGRES_DB) force $(version)

.PHONY: run build migrate-up migrate-down migrate-force
