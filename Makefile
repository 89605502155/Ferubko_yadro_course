.PHONY: all
all:
	go build -o myapp main.go
build: 
	go build -o myapp main.go
run:
	go run main.go