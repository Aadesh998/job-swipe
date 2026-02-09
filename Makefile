APP_NAME=job_swipe
CMD_PATH=./cmd/server/main.go
MIGRATION_PATH=./cmd/migration/main.go

build:
	go build -o $(APP_NAME).exe $(CMD_PATH)

run:
	go run $(CMD_PATH)

clean:
	del $(APP_NAME).exe

migration:
	go run $(MIGRATION_PATH)

run-gateway:
	go run ./cmd/gateway/main.go

compose-up:
	docker compose up -d

compose-down:
	docker compose down