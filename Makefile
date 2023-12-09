postgres: 
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

postgres-up: 
	docker start postgres

postgres-down: 
	docker stop postgres

createdb: 
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb: 
	docker exec -it postgres dropdb -U root simple_bank

migrate-up: 
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" --verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/simple_bank?sslmode=disable" --verbose down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test ./... -cover -v