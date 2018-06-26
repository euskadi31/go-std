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

travis:
	@go test -cover -covermode=count -coverprofile ./coverage.out ./...
