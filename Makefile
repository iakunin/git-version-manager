.PHONY: models

tag:
	git fetch --tags && \
	docker run \
	--tty \
	--interactive \
	--rm \
	--volume="${PWD}":/home \
	--workdir=/home \
	iakunin/git-version-manager:0.0.4 && \
	git push --tags

docker-build:
	VERSION=$$(git tag -l --sort=v:refname | tail -n 1); \
	docker build \
	--tag=iakunin/git-version-manager:$$VERSION \
	--tag=iakunin/git-version-manager:latest \
	.

docker-push:
	VERSION=$$(git tag -l --sort=v:refname | tail -n 1); \
	docker push iakunin/git-version-manager:$$VERSION && \
	docker push iakunin/git-version-manager:latest

deploy: tag docker-build docker-push
