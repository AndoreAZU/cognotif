initiate-database:
	docker-compose up -d

destroy-database:
	docker-compose down

run-core:
	go build -v -o core cmd/core/main.go && export APP_ENV=dev && ./core

generate-report:
	go run cmd/report/main.go

run-scheduler:
	go run cmd/scheduler/main.go