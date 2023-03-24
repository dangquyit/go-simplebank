makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
postgres:
	docker run --name postgres -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:latest
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable" -verbose down
sqlc:
	docker run --rm -v "$(makeFileDir):/src" -w /src kjconroy/sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/dangquyit/go-simplebank/db/sqlc Store
