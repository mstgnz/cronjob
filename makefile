include .env

# Makefile for building and running the Docker image and container
.PHONY: run live build db stop cleanI cleanC exec test
.DEFAULT_GOAL:= run

# Build the Docker image and run the container
run: cleanC build
	docker run -d --name $(APP_NAME) -p $(PORT):$(PORT) --env-file .env $(APP_NAME)

live:
#	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'make && docker logs --follow $(APP_NAME)'
	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'go build -o /tmp/build ./cmd && /tmp/build'

# Build the Docker image
build:
	docker build -t $(APP_NAME) .

# Postgres
db:
	docker volume create $(APP_NAME)-postgres
	docker run --name $(APP_NAME)-postgres \
		-p $(DB_PORT):5432  \
		-e POSTGRES_DB=$(DB_NAME) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-v $(APP_NAME)-postgres:/var/lib/postgresql/data -d \
		postgres

# Stop and remove the Docker container
# --time=600 for waiting running job
stop:
	docker stop --time=600 $(APP_NAME)
	docker rm $(APP_NAME)

# Run the application inside the Docker container
exec:
	docker exec -it $(APP_NAME) bash

# Clean up the Docker image
cleanI:
	docker rmi $(APP_NAME)
	docker builder prune --filter="image=$(APP_NAME)"

cleanC:
	@CONTAINER_EXISTS=$$(docker ps -aq --filter name=$(APP_NAME)); \
	if [ -n "$$CONTAINER_EXISTS" ]; then \
		echo "Stopping and removing container $(APP_NAME)"; \
		docker stop $(APP_NAME); \
		docker rm $(APP_NAME); \
	else \
		echo "No such container: $(APP_NAME)"; \
	fi

test: 
	go test -v ./...