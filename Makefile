run:
	go run cmd/main.go

build:
	go build -o theme ./cmd/main.go

install:
	go build -o ~/go/bin/theme ./cmd/main.go