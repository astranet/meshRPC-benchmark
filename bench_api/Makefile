all:

GIT_HASH := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMAGE_NAME := "docker.direct/meshrpc/benchmark/bench_api"
CURRDIR=$(PWD)

image:
	docker build --force-rm --build-arg GIT_HASH=$(GIT_HASH) -t $(IMAGE_NAME):local -f Dockerfile ..
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):$(GIT_HASH)
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):latest

push:
	docker push $(IMAGE_NAME):$(GIT_HASH)
	docker push $(IMAGE_NAME):latest

gen:
	go generate ./...

install:
	go get github.com/astranet/meshRPC-benchmark/bench_api

.PHONY: gen
