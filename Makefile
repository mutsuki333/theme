run:
	go run cmd/theme/main.go

build:
	go build -o theme ./cmd/theme/main.go

install:
	go build -o ~/go/bin/theme ./cmd/theme/main.go
