include .env

# Makefile for building and running the Docker image and container
.PHONY: run live build stop cleanI cleanC exec test
.DEFAULT_GOAL:= run

# Build the Docker image and run the container
run: cleanC build
	docker run -d --name $(PROJECT_NAME) -p $(PORT):$(PORT) --env-file .env $(PROJECT_NAME)

live:
	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'make && docker logs --follow $(PROJECT_NAME)'
#	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'go build -o /tmp/build ./cmd && /tmp/build'

# Build the Docker image
build:
	docker build -t $(PROJECT_NAME) .

# Stop and remove the Docker container
stop:
	docker stop --time=600 $(PROJECT_NAME)
	docker rm $(PROJECT_NAME)

# Run the application inside the Docker container
exec:
	docker exec -it $(PROJECT_NAME) bash

# Clean up the Docker image
cleanI:
	docker rmi $(PROJECT_NAME)
	docker builder prune --filter="image=$(PROJECT_NAME)"

cleanC:
	@CONTAINER_EXISTS=$$(docker ps -aq --filter name=$(PROJECT_NAME)); \
	if [ -n "$$CONTAINER_EXISTS" ]; then \
		echo "Stopping and removing container $(PROJECT_NAME)"; \
		docker stop $(PROJECT_NAME); \
		docker rm $(PROJECT_NAME); \
	else \
		echo "No such container: $(PROJECT_NAME)"; \
	fi

test: 
	go test -v ./...