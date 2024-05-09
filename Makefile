DEFOULT_STRING:=""
all:
	go build -c -i xkcd cmd/xkcd/main.go
build: 
	go build -c -i xkcd cmd/xkcd/main.go
run:
	go run cmd/xkcd/main.go -c "."
run1:
	go run cmd/xkcd/main.go -c "." -i -s "follower brings bunch of questions"
run2:
	go run cmd/xkcd/main.go -c "." -i -s "i'll follow you as long as you are following me"
runf:
	go run cmd/xkcd/main.go -c "." -i -s "$(DEFOULT_STRING)"