# Login to AWS registry (must have docker running)
docker-login:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 712639424220.dkr.ecr.us-west-2.amazonaws.com

# Build docker target
docker-build:
	docker build -f Dockerfile -t poly-webhook --force-rm --no-cache .

docker-tag:
	$(eval REV=$(shell git rev-parse HEAD | cut -c1-7))
	docker tag poly-webhook 712639424220.dkr.ecr.us-west-2.amazonaws.com/poly/asset:latest

# Push to registry
docker-push:
	docker push 712639424220.dkr.ecr.us-west-2.amazonaws.com/poly/asset:latest

# Build docker image and push to AWS registry
docker-build-and-push: docker-login docker-build docker-tag docker-push
