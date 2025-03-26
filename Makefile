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
  -e POSTGRES_DB=simple_bank \
  -v postgres-data:/var/lib/postgresql/data \
  postgres

pgadmin4:
	docker run -d --rm \
  --name pgadmin4 \
  --network bank-network \
  -p 8000:80 \
  -e PGADMIN_DEFAULT_EMAIL=admin@example.com \
  -e PGADMIN_DEFAULT_PASSWORD=adminSecret \
  -v pgadmin-data:/var/lib/pgadmin \
  dpage/pgadmin4

redis:
	docker run -d --rm --name redis -p 6379:6379 -d redis:7-alpine

mysql8up:
	docker run --rm --name mysql8 -p 3306:3306 \
	-e MYSQL_ROOT_PASSWORD=adminSecret \
	-v mysql8_data:/var/lib/mysql \
	-d mysql:8

mysql:
	docker ps | grep mysql8 || echo "Container is not running. Start it with 'make mysql8up'."
	docker exec -it mysql8 mysql -uroot -padminSecret -e "CREATE DATABASE IF NOT EXISTS simple_bank;"
	docker exec -it mysql8 mysql -uroot -padminSecret simple_bank

createdb:
ifdef DOCKER_EXEC
	docker exec -it postgres createdb --username=admin --owner=admin simple_bank
else
	psql -h postgres -U admin -d postgres -c "CREATE DATABASE simple_bank;"
endif

dropdb:
	docker exec -it postgres dropdb simple_bank

stopdb:
	docker stop postgres pgadmin4

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateUpVersion:
	migrate -path db/migration -database "$(DB_URL)" -verbose up $(version)

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateDownVersion:
	migrate -path db/migration -database "$(DB_URL)" -verbose down $(version)

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

frontend:
	cd frontend && npm run dev

docServer:
	docker run --rm --name simplebank -p 8080:8080 --network bank-network -e "DB_SOURCE=postgresql://admin:adminSecret@postgres:5432/simple_bank?sslmode=disable" simplebank

mock:
	~/go/bin/mockgen -package mockdb --destination db/mock/store.go simplebank/db/sqlc Store
	~/go/bin/mockgen -package mockwk --destination worker/mock/distributor.go simplebank/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

startTestEnv:
	@$(MAKE) postgres
	@sleep 2
	@$(MAKE) pgadmin4
	@sleep 2
	@$(MAKE) redis

.PHONY: startTestEnv network postgres mysql8up mysql createdb dropdb migrateup migratedown migrateUpVersion migrateDownVersion new_migration db_docs db_schema sqlcgen sqlcinit test server frontend docServer mock proto redis stopdb
