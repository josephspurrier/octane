# This Makefile is an easy way to run common operations.
#
# Tip: Each command is run on its own line so you can't CD unless you
# connect commands together using operators. See examples:
# A; B    # Run A and then B, regardless of success of A
# A && B  # Run B if and only if A succeeded
# A || B  # Run B if and only if A failed
# A &     # Run A in background.
# Source: https://askubuntu.com/a/539293
#
# Tip: Use $(shell app param) syntax when expanding a shell return value.

# Load the shared environment variables (can be shared with docker-compose.yml).
include .env

# Set local environment variables.
MYSQL_NAME=octane_db_1

.PHONY: run
run: swagger-gen  # Generate swagger and run.
	cd example/app/cmd/api && MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} go run main.go

.PHONY: swagger-get
swagger-get: # Download the Swagger generation tool.
	go get github.com/go-swagger/go-swagger/cmd/swagger

.PHONY: swagger-gen
swagger-gen: # Generate the Swagger spec.
	cd example/app/cmd/api && swagger generate spec -o ./swaggerui/swagger.json

.PHONY: db-init
db-init: # Launch database container.
	docker run -d --name=${MYSQL_NAME} -p 3306:3306 -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} ${MYSQL_CONTAINER}

.PHONY: db-start
db-start:
	# Start the stopped database container.
	docker start ${MYSQL_NAME}

.PHONY: db-stop
db-stop: # Stop the running database container.
	docker stop ${MYSQL_NAME}

.PHONY: db-reset
db-reset: # Drop the database, create the database, and perform the migrations.
	docker exec ${MYSQL_NAME} sh -c "exec mysql -h 127.0.0.1 -uroot -p${MYSQL_ROOT_PASSWORD} -e 'DROP DATABASE IF EXISTS main;'"
	docker exec ${MYSQL_NAME} sh -c "exec mysql -h 127.0.0.1 -uroot -p${MYSQL_ROOT_PASSWORD} -e 'CREATE DATABASE IF NOT EXISTS main DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;'"
	MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} go run ./example/app/cmd/dbmigrate/main.go

.PHONY: db-rm
db-rm: # Stop and remove the database container.
	docker rm -f ${MYSQL_NAME}