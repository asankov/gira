build:
	go build ./...

test:
	go test ./...

generate:
	go generate ./...

integration_tests:
	docker-compose -f cmd/integrationtests/docker-compose.yml up -d
	go get -u github.com/pressly/goose/cmd/goose
	# give the DB time to start
	sleep 5
	~/go/bin/goose -dir sql/ postgres 'host=localhost port=21665 user=gira dbname=gira password=password sslmode=disable' up
	go test cmd/integrationtests/*.go -tags integration_tests
