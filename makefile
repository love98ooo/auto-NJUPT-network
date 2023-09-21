APP_NAME = auto-NJUPT-network

windows:
	GOOS=windows GOARCH=amd64 go build -o out/${APP_NAME}.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o out/${APP_NAME}-linux-amd64
	GOARM=7 GOOS=linux GOARCH=arm go build -o out/${APP_NAME}-linux-armv7
    GOARM=8 GOOS=linux GOARCH=arm64 go build -o out/${APP_NAME}-linux-arm64
	GOOS=linux GOARCH=mips go build -o out/${APP_NAME}-linux-mips

darwin:
	GOOS=darwin GOARCH=amd64 go build -o out/${APP_NAME}-darwin-amd64


all: windows linux darwin

.PHONY: build all windows linux darwin
