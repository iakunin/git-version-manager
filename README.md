# git-version-manager

## Running
```shell
docker run \
    --tty \
    --interactive \
    --rm \
    --volume="${PWD}":/home \
    --workdir=/home \
    iakunin/git-version-manager:0.0.7 \
    /git-version-manager \
    --prefix=myPrefix \
    --suffix=myAwesomeSuffix
```
