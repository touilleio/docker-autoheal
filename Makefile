
VERSION=v1.0.0
GOOS=linux
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
DOCKER=docker
VERSION_MAJOR=$(shell echo $(VERSION) | cut -f1 -d.)
VERSION_MINOR=$(shell echo $(VERSION) | cut -f2 -d.)
BINARY_NAME=docker-autoheal
GO_PACKAGE=touille/docker-autoheal
DOCKER_REGISTRY=docker.touille.io
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

ensure:
	GOOS=${GOOS} $(GOCMD) mod vendor

clean:
	$(GOCLEAN)

lint:
	$(GOLINT) ...

build:
	GOOS=${GOOS} $(GOBUILD) \
		-ldflags "-X version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} \
				  -X version.BuildDate=${BUILD_DATE} \
				  -X version.Version=${VERSION}" \
		-o ${BINARY_NAME} .

package:
	$(DOCKER) build -f Dockerfile \
	  -t ${DOCKER_REGISTRY}/${GO_PACKAGE}:$(VERSION) \
	  -t ${DOCKER_REGISTRY}/${GO_PACKAGE}:$(VERSION_MAJOR).$(VERSION_MINOR) \
	  -t ${DOCKER_REGISTRY}/${GO_PACKAGE}:$(VERSION_MAJOR) \
	  .

test:
	$(GOTEST) ./...

release:
	$(DOCKER) push ${DOCKER_REGISTRY}/${GO_PACKAGE}:$(VERSION)
	$(DOCKER) push ${DOCKER_REGISTRY}/${GO_PACKAGE}:$(VERSION_MAJOR).$(VERSION_MINOR)
	$(DOCKER) push ${DOCKER_REGISTRY}/${GO_PACKAGE}:$(VERSION_MAJOR)
