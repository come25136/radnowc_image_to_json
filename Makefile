build:
	mkdir -p functions
	go get
	GOARCH=amd64 GOOS=linux go build -o ./functions/main *.go
