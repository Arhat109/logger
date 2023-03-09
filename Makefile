DOCKER_COMPOSE_PATH?=./deployments/docker_compose/docker-compose.test.yml

precommit:
		gofmt -w -s -d .
		golangci-lint run -v

start-unit-test:
	go test ./...
	exit $$test_status_code

install_linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint --version

run_linter:
	golangci-lint run -v

hook-install:
	cp ./pre-commit .git/hooks/

post-unpack: install-tools install_linter hook-install

