DB_URL=postgresql://admin:adminSecret@localhost:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run -d --rm \
  --name postgres \
  --network bank-network \
  -p 5432:5432 \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=adminSecret \
  -e POSTGRES_DB=simplebank \
  -v postgres-data:/var/lib/postgresql/data \
  postgres

pgadmin4:
	docker run -d --rm \
  --name pgadmin4 \
  --network bank-network \
  -p 8080:80 \
  -e PGADMIN_DEFAULT_EMAIL=admin@example.com \
  -e PGADMIN_DEFAULT_PASSWORD=adminSecret \
  -v pgadmin-data:/var/lib/pgadmin \
  dpage/pgadmin4

mysql:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=adminSecret -d mysql:8

createdb:
	docker exec -it postgres createdb --username=admin --owner=admin simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlcgen:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc generate

sqlcinit:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc init

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/techschool/simplebank/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlcgen sqlcinit test server mock proto evans redis
