# git-semver

See https://github.com/edgexfoundry/git-semver for the source of inspiration.

## Running
```shell
docker run \
    --tty \
    --interactive \
    --rm \
    --volume="${PWD}":/home \
    --workdir=/home \
    iakunin/git-semver:latest \
    /git-semver \
    --prefix=myPrefix \
    --suffix=myAwesomeSuffix
```
