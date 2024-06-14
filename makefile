include .env

# Makefile for building and running the Docker image and container
.PHONY: run live build db redis stop cleanI cleanC exec test
.DEFAULT_GOAL:= run

# Build the Docker image and run the container
run: cleanC create_network connect build
	docker run -d \
		--env-file .env \
		--restart always \
		--log-driver none \
		--name $(APP_NAME) \
		--network $(APP_NAME) \
		--port $(APP_PORT):$(APP_PORT) \
		$(APP_NAME)

live:
#	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'make && docker logs --follow $(APP_NAME)'
	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'go build -o /tmp/build ./cmd && /tmp/build'

# Build the Docker image
build:
	docker build -t $(APP_NAME) .

connect:
ifeq ($(APP_ENV),prod)
	$(MAKE) db
	$(MAKE) redis
else
	@echo "Skipping db target since APP_ENV is not 'prod'"
endif

# Postgres
db: create_volume
	docker run -d --name $(APP_NAME)-postgres \
		-p $(DB_PORT):5432  \
		--network $(APP_NAME) \
		-e POSTGRES_DB=$(DB_NAME) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-v $(APP_NAME):/var/lib/postgresql/data \
		postgres

# Redis
redis:
	docker run -d --name $(APP_NAME)-redis -p 6379:6379 --restart always --network $(APP_NAME) redis:latest

create_network:
	@if ! docker network inspect $(APP_NAME) >/dev/null 2>&1; then \
		docker network create $(APP_NAME); \
	else \
		echo "Network '$(APP_NAME)' already exists, using existing network."; \
	fi

create_volume:
	@if ! docker volume inspect $(APP_NAME) >/dev/null 2>&1; then \
		docker volume create $(APP_NAME); \
	else \
		echo "Volume '$(APP_NAME)' already exists, skipping creation."; \
	fi

# Stop and remove the Docker container
# --time=600 for waiting running job
stop:
	docker stop --time=600 $(APP_NAME)
	docker rm $(APP_NAME)

# Run the application inside the Docker container
exec:
	docker exec -it $(APP_NAME) $(CMD)

# Clean up the Docker image
cleanI:
	docker rmi $(APP_NAME)
	docker builder prune --filter="image=$(APP_NAME)"

cleanC:
	@CONTAINER_EXISTS=$$(docker ps -aq --filter name=$(APP_NAME)); \
	if [ -n "$$CONTAINER_EXISTS" ]; then \
		echo "Stopping and removing containers starting with $(APP_NAME)"; \
		CONTAINERS=$$(docker ps -aq --filter name=$(APP_NAME)); \
		for container in $$CONTAINERS; do \
			echo "Stopping and removing container $$container"; \
			docker stop $$container; \
			docker rm $$container; \
		done; \
	else \
		echo "No such containers starting with: $(APP_NAME)"; \
	fi

test: 
	go test -v ./...