.PHONY: models

tag:
	git fetch --tags && \
	docker run \
	--tty \
	--interactive \
	--rm \
	--volume="${PWD}":/home \
	--workdir=/home \
	iakunin/git-semver:0.0.1 \
	/git-semver && \
	git push --tags

docker-build:
	VERSION=$$(git tag -l --sort=v:refname | tail -n 1); \
	docker build \
	--tag=iakunin/git-semver:$$VERSION \
	--tag=iakunin/git-semver:latest \
	.

docker-push:
	docker push --all-tags iakunin/git-semver

deploy: tag docker-build docker-push
