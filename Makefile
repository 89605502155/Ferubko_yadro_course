all:
	go build -c xkcd cmd/xkcd/main.go
build: 
	go build -c xkcd cmd/xkcd/main.go
run:
	go run cmd/xkcd/main.go -c