# Makefile are useful when working in a team, 
# to easily setup the project on their local machine

# create new postgres container, change password as u wish
postgres: 
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mypass -d postgres:16

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres simple_bank

consoledb:
	docker exec -it postgres16 psql -U postgres simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:Anaana123@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:Anaana123@localhost:5433/simple_bank?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres16 dropdb simple_bank

.PHONY: postgres createdb consoledb migrateup migratedown dropdb

