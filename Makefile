DEFOULT_STRING:=""
NAME_MIGRATION:="inti"
all:
	go build -c -i xkcd cmd/xkcd/main.go
build: 
	go build -c -i xkcd cmd/xkcd/main.go
run:
	go run cmd/xkcd/main.go -c "."
runf:
	go run cmd/xkcd/main.go -c "$(DEFOULT_STRING)"

test_c:
	go test -cover ./...

test_v:
	go test -v ./...
test_r:
	go test -race ./...
test:
	go test ./...

n_migration:
	sql-migrate new sql $(NAME_MIGRATION)
mm:
	sql-migrate new sql up