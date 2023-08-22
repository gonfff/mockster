lint:
	golangci-lint run --config .golangci.yml

test:
	go test ./...

test-cover:
	go test -cover ./...

test-cover-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out