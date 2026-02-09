lint:
	golangci-lint run --config .golangci.yml

test:
	go test ./...

test-cover:
	go test -cover ./...

test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

race-test:
	go test -race -timeout=60s -count 1 ./...

run:
	go run app/main.go

build-app:
	go build -o ./mockster app/main.go

build-ui:
	npm run build:ui

check-ui:
	npm run check:ui

build:
	docker build -t $(tag) .
