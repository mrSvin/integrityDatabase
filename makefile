up:
	docker-compose up -d
down:
	docker-compose down

transfer test:
	go test -v ./integration-test -run Test_ServiceBatch
