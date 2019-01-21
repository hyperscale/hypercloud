BUILD_DIR ?= build
COMMIT = $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --always --tags --dirty)
ORG := github.com/hyperscale
PROJECT := hyperpaas
REPOPATH ?= $(ORG)/$(PROJECT)
VERSION_PACKAGE = $(REPOPATH)/pkg/hyperpaas/version

GO_LDFLAGS :="
GO_LDFLAGS += -X $(VERSION_PACKAGE).version=$(VERSION)
GO_LDFLAGS += -X $(VERSION_PACKAGE).buildDate=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitCommit=$(COMMIT)
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitTreeState=$(if $(shell git status --porcelain),dirty,clean)
GO_LDFLAGS +="

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: release
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

.PHONY: all
all: deps build test

.PHONY: deps
deps:
	@go mod download

.PHONY: clean
clean:
	@go clean -i ./...

generate: $(GO_FILES)
	@go generate ./...

$(BUILD_DIR)/coverage.out: $(GO_FILES)
	@CGO_ENABLED=0  go test -cover -coverprofile $(BUILD_DIR)/coverage.out.tmp ./...
	@cat $(BUILD_DIR)/coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > $(BUILD_DIR)/coverage.out
	@rm $(BUILD_DIR)/coverage.out.tmp

.PHONY: ci-test
ci-test:
	@go test -race -cover -coverprofile ./coverage.out.tmp -v ./... | go2xunit -fail -output tests.xml
	@cat ./coverage.out.tmp | grep -v '.pb.go' | grep -v 'mock_' > ./coverage.out
	@rm ./coverage.out.tmp
	@echo ""
	@go tool cover -func ./coverage.out

.PHONY: lint
lint:
	@CGO_ENABLED=0 golangci-lint run ./...

.PHONY: test
test: $(BUILD_DIR)/coverage.out

.PHONY: coverage
coverage: $(BUILD_DIR)/coverage.out
	@echo ""
	@go tool cover -func ./$(BUILD_DIR)/coverage.out

.PHONY: coverage-html
coverage-html: $(BUILD_DIR)/coverage.out
	@go tool cover -html ./$(BUILD_DIR)/coverage.out


# Build targets

${BUILD_DIR}/hyperpaas-starter: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

${BUILD_DIR}/hyperpaas-installer: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

${BUILD_DIR}/hyperpaas-server: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

.PHONY: build
build: ${BUILD_DIR}/hyperpaas-starter ${BUILD_DIR}/hyperpaas-installer ${BUILD_DIR}/hyperpaas-server


# Docker targets

docker: docker-hyperpaas-starter docker-hyperpaas-installer docker-hyperpaas-server docker-hyperpaas-manager

.PHONY: docker-hyperpaas-starter
docker-hyperpaas-starter: _docker-hyperpaas-starter

.PHONY: docker-hyperpaas-installer
docker-hyperpaas-installer: _docker-hyperpaas-installer

.PHONY: docker-hyperpaas-server
docker-hyperpaas-server: _docker-hyperpaas-server

.PHONY: docker-hyperpaas-manager
docker-hyperpaas-manager: _docker-hyperpaas-manager

_docker-%:
	@docker build -f cmd/$*/Dockerfile -t 127.0.0.1:5000/$*:latest .
	@docker image push 127.0.0.1:5000/$*


# Run targets

run: docker
	@sudo docker run -p 8578:8080 \
		-e "USERNAME=dacteev" \
		-e "PASSWORD=test" \
		-v $(shell pwd)/var/lib/hyperpaas:/var/lib/hyperpaas \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--rm $(IMAGE)

.PHONY: run-hyperpaas-server
run-hyperpaas-server: ${BUILD_DIR}/hyperpaas-server
	@echo "Running $<..."
	@./$< --config=./cmd/$*/config.yml

.PHONY: run-hyperpaas-starter
run-hyperpaas-starter: ${BUILD_DIR}/hyperpaas-starter
	@echo "Running $<..."
	@./$<


# Swarm targets

.PHONY: stack-deploy-dev
stack-deploy-dev:
	@docker stack deploy -c dev/docker-compose.yml acme

.PHONY: stack-deploy-installer
stack-deploy-installer:
	@docker stack deploy -c cmd/hyperpaas-installer/docker-compose.yml hyperpaas
