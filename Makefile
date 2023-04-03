makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
postgres:
	docker run --name postgres --network bank-network -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:latest
createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5431/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	docker run --rm -v "$(makeFileDir):/src" -w /src kjconroy/sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/dangquyit/go-simplebank/db/sqlc Store
