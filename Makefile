# Phony targets
.PHONY: build test docker-build docker-run docker-push

IMAGE_NAME=ghcr.io/xonvanetta/kubernetes-git-sync

# Build the Go binary
build:
	CGO_ENABLED=0 GOOS=linux go build ./cmd/kubernetes-git-sync

# Run tests
test:
	go test ./...

# Build the Docker image
docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) -f build/Dockerfile .

# Run the Docker container
docker-run:
	docker run --rm --name kubernetes-git-sync -u 1000 -v /home/$(USER)/.kube:/.kube -v $(PWD)/build/.gitconfig:/.gitconfig $(IMAGE_NAME):$(TAG)

# Push the Docker image to the registry
docker-push:
	docker tag $(IMAGE_NAME):$(TAG) $(REGISTRY)/$(IMAGE_NAME):$(TAG)
	docker push $(REGISTRY)/$(IMAGE_NAME):$(TAG)