include .env 
export

APP_NAME = invoGenius
SRC_DIR = .

all: run

build:
	go build -o bin/$(APP_NAME) $(SRC_DIR)

run: air

test:
	go test ./...

clean:
	rm -rf bin/
	go clean

create_migrate:
	@if [ "$(word 2, $(MAKECMDGOALS))" = "" ]; then \
		echo "Error: please provide a migration name. Usage: make create_migrate <name>"; \
		exit 1; \
	fi
	migrate create -seq -ext sql -dir ./db/migrations $(word 2, $(MAKECMDGOALS))

# This is required to avoid errors when Make sees the second word as a target
%:
	@:

migrate_up:
	migrate -path=./db/migrations -database="mysql://$(DATABASE_URI)" up

migrate_down:
	migrate -path=./db/migrations -database="mysql://$(DATABASE_URI)" down

.PHONY: all build run test clean create_migrate migrate_up migrate_down
