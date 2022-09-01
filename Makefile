BINARY=scrum-bot

default: build run clean

build:
	go mod tidy
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY}-darwin main.go
	@# GOARCH=amd64 GOOS=linux go build -o ${BINARY}-linux main.go

run:
	echo ""; ./${BINARY}-darwin config.json

clean:
	@go clean
	@rm ${BINARY}-darwin
	@# rm ${BINARY}-linux

predeploy:
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o ${BINARY}-linux main.go
