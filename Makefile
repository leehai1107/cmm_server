run: 
	go run main.go api
lint:
	golangci-lint run
docker-up:
	docker compose up 
swagger:
	swag init