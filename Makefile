build:
	go build ./...

test:
	go test ./...

generate:
	go generate ./...

integration_tests:
	go test -v cmd/integrationtests/*.go -tags integration_tests
