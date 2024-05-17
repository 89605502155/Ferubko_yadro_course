DEFOULT_STRING:=""
NAME_MIGRATION:="inti"
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

n_migration:
	sql-migrate new sql $(NAME_MIGRATION)
mm:
	sql-migrate new sql up