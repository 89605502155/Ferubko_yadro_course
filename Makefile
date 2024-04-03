DEFOULT_N:=5
.PHONY: all
all:
	go build -o xkcd cmd/xkcd/main.go
build: 
	go build -o xkcd cmd/xkcd/main.go
run:
	go run cmd/xkcd/main.go

run_o:
	go run cmd/xkcd/main.go -o
run_n:
	go run cmd/xkcd/main.go -n ${DEFOULT_N}