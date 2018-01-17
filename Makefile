install:
	go install -v

build:
	go build -v

build_pi:
	GOOS=linux GOARCH=arm GOARM=6 go build -v

test:
	go test --race -v ./...

.PHONY: install build build_pi test