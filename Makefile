BUILD_DIR ?= build
COMMIT = $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --always --tags --dirty)
ORG := github.com/hyperscale
PROJECT := hypercloud
REPOPATH ?= $(ORG)/$(PROJECT)
VERSION_PACKAGE = $(REPOPATH)/pkg/hypercloud/version

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

create-build-dir:
	@mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/coverage.out: create-build-dir $(GO_FILES) go.mod go.sum
	@go test -race  -cover -coverprofile $(BUILD_DIR)/coverage.out.tmp ./...
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

${BUILD_DIR}/hypercloud-starter: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

${BUILD_DIR}/hypercloud-installer: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

${BUILD_DIR}/hypercloud-server: $(GO_FILES)
	@echo "Building $@..."
	@go generate ./cmd/$(subst ${BUILD_DIR}/,,$@)/
	@go build -ldflags $(GO_LDFLAGS) -o $@ ./cmd/$(subst ${BUILD_DIR}/,,$@)/

.PHONY: build
build: ${BUILD_DIR}/hypercloud-starter ${BUILD_DIR}/hypercloud-installer ${BUILD_DIR}/hypercloud-server


# Docker targets

docker: docker-hypercloud-starter docker-hypercloud-installer docker-hypercloud-server docker-hypercloud-manager

.PHONY: docker-hypercloud-starter
docker-hypercloud-starter: _docker-hypercloud-starter

.PHONY: docker-hypercloud-installer
docker-hypercloud-installer: _docker-hypercloud-installer

.PHONY: docker-hypercloud-server
docker-hypercloud-server: _docker-hypercloud-server

.PHONY: docker-hypercloud-manager
docker-hypercloud-manager: _docker-hypercloud-manager

_docker-%:
	@docker build -f cmd/$*/Dockerfile -t 127.0.0.1:5000/$*:latest .
	@docker image push 127.0.0.1:5000/$*


# Run targets

run: docker
	@sudo docker run -p 8578:8080 \
		-e "USERNAME=dacteev" \
		-e "PASSWORD=test" \
		-v $(shell pwd)/var/lib/hypercloud:/var/lib/hypercloud \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--rm $(IMAGE)

.PHONY: run-hypercloud-server
run-hypercloud-server: ${BUILD_DIR}/hypercloud-server
	@echo "Running $<..."
	./$< --config=./cmd/$(subst ${BUILD_DIR}/,,$<)/config.yml

.PHONY: run-hypercloud-starter
run-hypercloud-starter: ${BUILD_DIR}/hypercloud-starter
	@echo "Running $<..."
	@./$<


# Swarm targets

.PHONY: stack-deploy-dev
stack-deploy-dev:
	@docker stack deploy -c dev/docker-compose.yml acme

.PHONY: stack-deploy-installer
stack-deploy-installer:
	@docker stack deploy -c cmd/hypercloud-installer/docker-compose.yml hypercloud
