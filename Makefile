.PHONY: all clean fmt vet test

PACKAGES = $(shell go list ./...)

release:
	@echo "Release v$(version)"
	@git pull
	@git checkout master
	@git pull
	@git checkout develop
	@git flow release start $(version)
	@git flow release finish $(version) -p -m "Release v$(version)"
	@git checkout develop
	@echo "Release v$(version) finished."

all: test

clean:
	@go clean -i ./...

test:
	@go test -cover -coverprofile ./coverage.out ./...

cover: test
	@echo ""
	@go tool cover -func ./coverage.out


.PHONY: lint
lint:
ifeq (, $(shell which golangci-lint))
	@echo "Install golangci-lint..."
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.46.2
endif
	@echo "lint..."
	@golangci-lint run --timeout=300s ./...
