run:
	docker compose up -d
	go run main.go

clean:
	docker compose down
