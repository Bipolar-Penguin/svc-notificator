build-push-docker:
	docker build --pull . -t kuwerin/svc-notificator:latest
	docker push kuwerin/svc-notificator:latest

check:
	golangci-lint run -c .golangci.yml

fmt:
	go fmt ./...

test:
	go test -v ./...

run-dev:
	docker-compose -f docker/docker-compose.yml -p svc-notificator up --build
