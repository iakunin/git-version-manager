
build-and-push-docker-image: docker-build docker-push

docker-build:
	docker build -t iakunin/git-semver:latest .

docker-push:
	docker push iakunin/git-semver:latest
