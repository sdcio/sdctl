# container image
REMOTE_REGISTRY := registry.kmrd.dev/sdcio/sdctl
TAG := $(shell git describe --tags)
IMAGE := $(REMOTE_REGISTRY):$(TAG)
USERID := 10000
# go versions
TARGET_GO_VERSION := go1.21.4
GO_FALLBACK := go
# We prefer $TARGET_GO_VERSION if it is not available we go with whatever go we find ($GO_FALLBACK)
GO_BIN := $(shell if [ "$$(which $(TARGET_GO_VERSION))" != "" ]; then echo $$(which $(TARGET_GO_VERSION)); else echo $$(which $(GO_FALLBACK)); fi)
USERID := 10000

build:
	mkdir -p bin
	CGO_ENABLED=0 ${GO_BIN} build -o bin/sdctl main.go

docker-build:
	ssh-add ./keys/id_rsa 2>/dev/null; true
	docker build --build-arg USERID=$(USERID) . -t $(IMAGE) --ssh default=$(SSH_AUTH_SOCK)

