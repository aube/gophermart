.PHONY: auth
auth:
	go run -v ./cmd/auth

.PHONY: mocks
mocks:
	mockgen -source=./internal/api/interfaces.go -destination=./mocks/api/mocks.go

.PHONY: test
test:
	go test -v -timeout 30s ./...

.PHONY: race
race:
	go test -v -race -timeout 30s ./...

.PHONY: cover
cover:
	go test -v -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o test-coverage.html
	rm coverage.out

#.DEFAULT_GOAL := auth
