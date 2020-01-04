# Only edit this file if it's in the base folder
WORKSHOP := $(shell basename "$$PWD")
PACKAGE_NAME ?= workshops/$(WORKSHOP)
CONTAINER_DIR := /go/src/$(PACKAGE_NAME)
GO_VERSION := 1.9.5
GOARCH := amd64
UNAME_S ?= $(shell uname -s)

.PHONY: default
default: build

init:
	which dep || go get -u github.com/golang/dep/cmd/dep
	dep ensure -v -vendor-only

build:
ifeq ($(UNAME_S),Linux)
	go build -ldflags "-X main.version=${WORKSHOP}" workshops/${WORKSHOP}/cmd/hello
endif
ifeq ($(UNAME_S),Darwin)
	GOOS=darwin go build -ldflags "-X main.version=${WORKSHOP}" -o hello_darwin workshops/${WORKSHOP}/cmd/hello
endif

run: build
	./hello

install:
	CGO_ENABLED=0 go install -ldflags "-X main.version=${WORKSHOP}" workshops/${WORKSHOP}/cmd/hello

image: PROJECT_ID = $(shell gcloud config get-value core/project)
image:
	docker build --build-arg workshop=${WORKSHOP} -t gcr.io/${PROJECT_ID}/hello:${WORKSHOP} .

publish: PROJECT_ID = $(shell gcloud config get-value core/project)
publish:
	docker push gcr.io/${PROJECT_ID}/hello:${WORKSHOP}

docker_%: PROJECT_ID = $(shell gcloud config get-value core/project)
docker_%:
	# create a container
	$(eval CONTAINER_ID := $(shell docker create \
		-i \
		-w $(CONTAINER_DIR) \
		-e PROJECT_ID=$(PROJECT_ID) \
		-e UNAME_S=$(UNAME_S) \
		golang:${GO_VERSION} \
		/bin/bash -c "make $*" \
	))

	# copy stuff into container
	find . | tar cf -  -T - | docker cp - $(CONTAINER_ID):$(CONTAINER_DIR)

	# run task inside container
	docker start -a -i $(CONTAINER_ID)

	# copy out any artifacts if any
ifeq ($(UNAME_S),Linux)
	docker cp $(CONTAINER_ID):$(CONTAINER_DIR)/hello hello || true
endif
ifeq ($(UNAME_S),Darwin)
	docker cp $(CONTAINER_ID):$(CONTAINER_DIR)/hello_darwin hello || true
endif
	# remove container
	docker rm $(CONTAINER_ID)
