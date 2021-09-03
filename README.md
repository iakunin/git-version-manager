# git-version-manager

## Running
```shell
docker run \
    --tty \
    --interactive \
    --rm \
    --volume="${PWD}":/home \
    --workdir=/home \
    iakunin/git-semver:0.0.4 \
    /git-semver \
    --prefix=myPrefix \
    --suffix=myAwesomeSuffix
```
