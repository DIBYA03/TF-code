DB_USER=wise_test
DB_PASSWORD=bypass
DB_NAME=wise
DB_WRITE_PORT=5432
DB_WRITE_ENDPOINT=localhost
DB_SSL_MODE=disable

build:
	go build -o goose *.go

up:
	docker run --rm  --name wisedb -p $(DB_WRITE_PORT):5432 -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_DB=$(DB_NAME) -d postgres:10.8
	docker ps -a

down:
	docker stop wisedb

prompt:
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_WRITE_ENDPOINT) -p $(DB_WRITE_PORT) -U $(DB_USER) -d $(DB_NAME)

migrate: build
	# Call this like: $ make args="up-to 123" migrate
	./goose postgres "dbname=$(DB_NAME) user=$(DB_USER) password=$(DB_PASSWORD) host=$(DB_WRITE_ENDPOINT) sslmode=$(DB_SSL_MODE)" $(if $(args),$(args), up)

.PHONY: docs
docs:
	mkdir -p ./docs
	rm -f ./docs/schemaspy.properties
	echo "schemaspy.t=pgsql" >> ./docs/schemaspy.properties
	echo "schemaspy.host=wisedb" >> ./docs/schemaspy.properties
	echo "schemaspy.port=$(DB_WRITE_PORT)" >> ./docs/schemaspy.properties
	echo "schemaspy.db=$(DB_NAME)" >> ./docs/schemaspy.properties
	echo "schemaspy.u=$(DB_USER)" >> ./docs/schemaspy.properties
	echo "schemaspy.p=$(DB_PASSWORD)" >> ./docs/schemaspy.properties
	docker run --link wisedb --rm -v "$(PWD)/docs:/output" -v "$(PWD)/docs/schemaspy.properties:/schemaspy.properties" schemaspy/schemaspy:latest
