postgres:
	docker run --name mydb --network bank-network -p 5432:5432  -e POSTGRES_USER=kadera -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it mydb createdb --username=kadera --owner=kadera simple_bank
dropdb:
	docker exec -it mydb dropdb --username=kadera simple_bank
migrate-up:
	migrate -path db/migration -database "postgresql://kadera:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrate-up1:
	migrate -path db/migration -database "postgresql://kadera:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migrate-down:
	migrate -path db/migration -database "postgresql://kadera:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migrate-down1:
	migrate -path db/migration -database "postgresql://kadera:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Yogksai/simplebank/db/sqlc Store
.PHONY: 
	createdb dropdb postgres migrate-up migrate-down sqlc test server mock migrate-down1 migrate-up1