BINARY_NAME=myapp
DSN="host=127.0.0.1 port=5432 user=postgres password=password dbname=concurrency sslmode=disable"
REDIS="127.0.0.1:6379"

## build: Buil binary
build:
	@echo "Building ..."
	env CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/web
	@echo "Built !"

## run: build and runs the application
run: build
	@echo "Starting..."
	env DSN=${DSN} REDIS=${REDIS} ./${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping ..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the applicaiton
restart: stop start

## test: run all tests
test:
	go test -v ./...