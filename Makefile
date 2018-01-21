.PHONY: all clean deps fmt vet test docker

EXECUTABLE ?= hyperpaas
IMAGE ?= hyperscale/$(EXECUTABLE)
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='-dev' --always)
COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X "hyperpaas.Revision=$(COMMIT)" -X "hyperpaas.Version=$(VERSION)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

release:
	@echo "Release v$(version)"
	@git pull
	@git checkout master
	@git pull
	@git checkout develop
	@git flow release start $(version)
	@echo "$(version)" > .version
	@sed -e "s/version: .*/version: \"v$(version)\"/g" docs/swagger.yaml > docs/swagger.yaml.new && rm -rf docs/swagger.yaml && mv docs/swagger.yaml.new docs/swagger.yaml
	@git add .version docs/swagger.yaml
	@git commit -m "feat(project): update version file" .version docs/swagger.yaml
	@git flow release finish $(version)
	@git push
	@git push --tags
	@git checkout master
	@git push
	@git checkout develop
	@echo "Release v$(version) finished."

all: deps build test

clean:
	@go clean -i ./...

deps:
	@glide install

fmt:
	@go fmt $(PACKAGES)

vet:
	@go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -ldflags '-s -w $(LDFLAGS)' -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

travis:
	@for PKG in $(PACKAGES); do go test -ldflags '-s -w $(LDFLAGS)' -cover -covermode=count -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

cover: test
	@echo ""
	@for PKG in $(PACKAGES); do go tool cover -func $$GOPATH/src/$$PKG/coverage.out; echo ""; done;

docker:
	#@sudo docker build --no-cache=true --rm -t $(IMAGE) .
	@sudo docker build --rm -t $(IMAGE) .

publish: docker
	@sudo docker tag $(IMAGE) $(IMAGE):latest
	@sudo docker push $(IMAGE)

bindata.go: docs/index.html docs/swagger.yaml
	@echo "Bin data..."
	@go-bindata docs/

$(EXECUTABLE): $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building $(EXECUTABLE)..."
	@CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)' ./cmd/hyperpaas/

build: $(EXECUTABLE)

run: docker
	@sudo docker run -p 8578:8080 \
		-e "USERNAME=dacteev" \
		-e "PASSWORD=test" \
		-v $(shell pwd)/var/lib/hyperpaas:/var/lib/hyperpaas \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--rm $(IMAGE)

docker-dev:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)' ./cmd/hyperpaas/
	@sudo docker build --rm -t $(IMAGE) -f Dockerfile.dev .

dev: docker-dev
	@sudo docker run --rm -p 8181:8080 \
		-e "USERNAME=dacteev" \
		-e "PASSWORD=test" \
		-v $(shell pwd)/var/lib/hyperpaas:/var/lib/hyperpaas \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(shell pwd)/ui/dist:/opt/hyperpaas/ui \
		--rm $(IMAGE)

dev-ui:
	@cd ui; ng build --sourcemaps --watch --base-href=/ui/ --aot -dev

build-hyperpaas-server: $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building hyperpaas-server..."
	@go generate ./cmd/hyperpaas-server/
	@CGO_ENABLED=0 go build ./cmd/hyperpaas-server/

run-hyperpaas-server: build-hyperpaas-server
	./hyperpaas-server

build-hyperpaas-starter: $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building hyperpaas-starter..."
	@go generate ./cmd/hyperpaas-starter/
	@CGO_ENABLED=0 go build ./cmd/hyperpaas-starter/

run-hyperpaas-starter: build-hyperpaas-starter
	./hyperpaas-starter

build-hyperpaas-installer: $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building hyperpaas-installer..."
	@go generate ./cmd/hyperpaas-installer/
	@CGO_ENABLED=0 go build ./cmd/hyperpaas-installer/

run-hyperpaas-installer: build-hyperpaas-installer
	./hyperpaas-installer

stack-deploy-dev:
	@docker stack deploy -c dev/docker-compose.yml acme

stack-deploy-installer:
	@docker stack deploy -c cmd/hyperpaas-installer/docker-compose.yml hyperpaas
