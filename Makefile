export POSTGRES_DSN=postgresql://postgres:postgres@localhost:5432/mydb?sslmode=disable&timezone=UTC
export MEMCACHED=localhost:11211

# GO
# ----------------------------------------

go-run:
	go build -o ./bin/poll-service ./cmd/. && ./bin/poll-service

go-build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/poll-service ./cmd/.

# DOCKER
# ----------------------------------------

docker-build: go-build
	docker build --tag dmitrukv/poll-service:1.0.0 -f ./docker/Dockerfile.dev .

docker-run:
	docker run --rm -p 8080:8080 dmitrukv/poll-service:1.0.0
  
# DOCKER-COMPOSE
# ----------------------------------------

dc-build-all: go-build
	docker-compose up -d --build

dc-build-backend: go-build
	docker-compose rm -svf backend
	docker-compose up -d --build backend

dc-up:
	docker-compose up -d

dc-up-backend:
	docker-compose rm -svf backend
	docker-compose up -d backend
dc-up-postgres:
	docker-compose rm -svf postgres
	docker-compose up -d postgres
dc-up-memcached:
	docker-compose rm -svf memcached
	docker-compose up -d memcached

dc-down:
	docker-compose down

dc-ps:
	docker-compose ps --all

# DOCKER-COMPOSE LOCAL
# ----------------------------------------

dc-local-ps:
	docker-compose -f ./docker-compose.local.yaml ps --all

dc-local-up:
	docker-compose -f ./docker-compose.local.yaml  up -d

dc-local-down:
	docker-compose -f ./docker-compose.local.yaml down