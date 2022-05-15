# CGO_ENABLED=0 == static by default
GO		?= go
GOOS		?= linux
DOCKER		?= docker
# -s removes symbol table and -ldflags -w debugging symbols
LDFLAGS		?= -asmflags -trimpath -ldflags "-s -w"
GOARCH		?= amd64
BINARY		?= pasto
CGO_ENABLED	?= 1

.PHONY: all analysis test

build: test lint
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	     $(GO) build $(LDFLAGS) -o build/$(BINARY)_$(GOOS)_$(GOARCH) cmd/pasto/pasto.go

lint:
	golangci-lint run ./...

docker-build:
	$(DOCKER) build --rm --target app -t $(BINARY)-test .

docker-run:
	docker run --rm -v $(shell pwd):/app/config $(BINARY)-test &

test:
	go test ./...

clean:
	-rm -rf build/

