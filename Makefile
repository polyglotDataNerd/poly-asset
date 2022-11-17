# Login to AWS registry (must have docker running)
docker-login:
	aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 426063958398.dkr.ecr.us-east-1.amazonaws.com

# Build docker target
docker-build:
	docker build -f Dockerfile -t curio-webhook --force-rm --no-cache .

docker-tag:
	$(eval REV=$(shell git rev-parse HEAD | cut -c1-7))
	docker tag curio-webhook 426063958398.dkr.ecr.us-east-1.amazonaws.com/curio/asset:latest

# Push to registry
docker-push:
	docker push 426063958398.dkr.ecr.us-east-1.amazonaws.com/curio/asset:latest

# Build docker image and push to AWS registry
docker-build-and-push: docker-login docker-build docker-tag docker-push
