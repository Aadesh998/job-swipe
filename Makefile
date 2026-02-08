APP_NAME=aron_project
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
