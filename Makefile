# Create db
createdb:
	 psql -h localhost -U postgres \
        -c "CREATE DATABASE activator WITH OWNER = postgres ENCODING = 'UTF8'";

# Drop db
dropdb:
	psql -h localhost -U postgres \
        -c "DROP DATABASE IF EXISTS activator";

migrateup:
	psql -h localhost -U postgres -d activator -a -f db/create_tables.sql

migratedown:
	psql -h localhost -U postgres -d activator -a -f db/delete_tables.sql

build:
	go build -o bin/app -v ./cmd/main.go

run:
	echo "starting service ..."
	bin/app

.DEFAULT_GOAL := build
