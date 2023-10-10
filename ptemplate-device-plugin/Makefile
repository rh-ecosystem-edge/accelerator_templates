EXEC_NAME := ptemplate-device-plugin
BASE_IMAGE ?= scratch
GO_FILES ?= $$(find . -name '*.go' -not -path './vendor/*')
REPO ?= quay.io/chrisp262/pt-device-plugin
TAG ?= plugin-latest
IMG ?= $(REPO):$(TAG)

CONSUMER_REPO=quay.io/chrisp262/pt-device-plugin
CONSUMER_TAG=consumer-latest
CONSUMER_IMG=$(CONSUMER_REPO):$(CONSUMER_TAG)

CONTAINER_RUNTIME_COMMAND ?= docker

build:
	go build -o $(EXEC_NAME) ./main.go

fmt:
	go fmt ./...

.PHONY: image
image: build
	$(CONTAINER_RUNTIME_COMMAND) build --build-arg BASE_IMAGE=$(BASE_IMAGE) -t $(IMG) .

.PHONY: image-push
image-push:
	$(CONTAINER_RUNTIME_COMMAND) push $(IMG)

consumer: 
	$(CONTAINER_RUNTIME_COMMAND) build -f Dockerfile.consumer -t $(CONSUMER_IMG) .

consumer-push:
	$(CONTAINER_RUNTIME_COMMAND) push $(CONSUMER_IMG)

clean:
	rm -f $(EXEC_NAME)

#.PHONY: build image push clean
