
PROJECT_ID ?= $(shell yc config get cloud-id)

CURDIR = $(shell pwd)
TAG = $(shell git rev-parse --short HEAD)
APP_NAME = $(shell basename $(CURDIR))

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ DOCKER Deployment
docker-build:  ## Build the docker image locally. Example of running : make docker-build PROJECT_ID=crp0n3m1vtkc60qrmi77
	docker build -t cr.yandex/$(PROJECT_ID)/$(APP_NAME):$(TAG) -f ./deployment/docker/catgpt/Dockerfile .

docker-push: ## Push the docker image to the registry
	docker push -t cr.yandex/$(PROJECT_ID)/$(APP_NAME):$(TAG) .