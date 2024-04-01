DEFOULT_STRING:=""
.PHONY: all
all:
	go build -o myapp main.go
build: 
	go build -o myapp main.go
run:
	go run main.go
run_with_first_arg:
	go run main.go -s "follower brings bunch of questions"
run_with_second_arg:
	go run main.go -s "i'll follow you as long as you are following me"
run_with_s_flag:
	go run main.go -s "$(DEFOULT_STRING)"
