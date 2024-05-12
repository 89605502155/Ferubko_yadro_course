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
	MIGRATE_FILE=migrate.toml migrate create -dir ./migrations -ext sql $(NAME_MIGRATION)
mm:
	migrate -path migrations -database "sqlite3://./xkcd.db?_auth&_auth_user=admin&_auth_pass=qwerty&_auth_crypt=sha256&_auth_salt=yadro&journal_mode=wal&cache=shared&mode=rwc" up