include .env

PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

.PHONY: run live build db redis stop cleanI cleanC exec test help
.DEFAULT_GOAL:= run

help: makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## run: Build the Docker image and run the container
run: cleanC create_network connect build
	docker run -d \
		--env-file .env \
		--restart always \
		--name $(APP_NAME) \
		--network $(APP_NAME) \
		-p $(APP_PORT):$(APP_PORT) \
		$(APP_NAME)

## live: Go build and running
live:
#	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'make && docker logs --follow $(APP_NAME)'
	find . -type f \( -name '*.go' -o -name '*.gohtml' \) | entr -r sh -c 'go build -o /tmp/build ./cmd && /tmp/build'

## build: Build the Docker image
build:
	docker build -t $(APP_NAME) .

connect:
ifeq ($(APP_ENV),prod)
	$(MAKE) db
	$(MAKE) redis
else
	$(MAKE) db
	@echo "Skipping db target since APP_ENV is not 'prod'"
endif

## db: Run Postgres
db: create_volume
	docker run -d --name $(APP_NAME)-postgres \
		-p $(DB_PORT):5432  \
		--network $(APP_NAME) \
		-e POSTGRES_DB=$(DB_NAME) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-v $(APP_NAME):/var/lib/postgresql/data \
		postgres

## redis: Run Redis
redis:
	docker run -d --name $(APP_NAME)-redis -p 6379:6379 --restart always --network $(APP_NAME) redis:latest

## create_network: Create network for this project name
create_network:
	@if ! docker network inspect $(APP_NAME) >/dev/null 2>&1; then \
		docker network create $(APP_NAME); \
	else \
		echo "Network '$(APP_NAME)' already exists, using existing network."; \
	fi

## create_volume: Create volume for this project name
create_volume:
	@if ! docker volume inspect $(APP_NAME) >/dev/null 2>&1; then \
		docker volume create $(APP_NAME); \
	else \
		echo "Volume '$(APP_NAME)' already exists, skipping creation."; \
	fi

## stop: Stop and remove the Docker container 
stop:
	docker stop --time=600 $(APP_NAME)
	docker rm $(APP_NAME)

## exec: Run the application inside the Docker container
exec:
	docker exec -it $(APP_NAME) $(CMD)

## cleanI: Clean up the Docker image
cleanI:
	docker rmi $(APP_NAME)
	docker builder prune --filter="image=$(APP_NAME)"

## cleanC: Clean up the Docker containers
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

## test: Run all test
test: 
	go test -v ./...