go-run:
	go build -o ./bin/polls-service ./cmd/. && ./bin/polls-service

go-build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/polls-service ./cmd/.

dc-build:
	docker-compose -f ./infrastructure/docker-compose.yaml up -d --build

dc-up:
	docker-compose -f ./infrastructure/docker-compose.yaml up -d

dc-down:
	docker-compose -f ./infrastructure/docker-compose.yaml down

dc-ps:
	docker-compose -f ./infrastructure/docker-compose.yaml ps -a