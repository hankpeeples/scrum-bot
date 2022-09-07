BINARY=scrum-bot
RASPBERRY_PI_ARM_VERSION=7

default: build run clean

build_darwin:
	go mod tidy
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY}-darwin main.go

build_linux:
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o ${BINARY}-linux main.go

build_raspberry-pi:
	go mod tidy
	GOARCH=arm GOARM=${RASPBERRY_PI_ARM_VERSION} GOOS=linux go build -o ${BINARY}-raspberry-pi main.go

run_darwin:
	echo "darwin"; ./${BINARY}-darwin

run_linux:
	echo "linux"; ./${BINARY}-linux

run_raspberry-pi:
	echo "raspberry-pi"; ./${BINARY}-raspberry-pi

clean:
	@go clean
	@rm ${BINARY}-darwin
	@rm ${BINARY}-linux
	@rm ${BINARY}-raspberry-pi
