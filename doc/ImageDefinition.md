# Image definition

Packager uses Dockerfiles for Image definitions. The definitions of images must be in Package
Context in `docker` directory.

## Requirements

The image definitions must comply with [specified requirements](./DockerContainerRequirements.md).

## Build context

Docker build context is a set of files which can be accessed from Dockerfile. Packager sets build
context as a directory where each Dockerfile is present in Package Context. So by default no other
files (etc in parent directories) can not be accessed. For a use case, when user wants to add some
common files for multiple image definitions, Packager adds additional build context which is set
to Package Context root directory.

For example to access `<package_context_root>/config/config.txt`, following Dockerfile code can be
used:

```dockerfile
# Copy from additional build context set to Package Context root
COPY --from=package-context config/config.txt /etc

# Without specifying "--from" only files next to Dockerfile can be accessed
COPY dir_next_to_dockerfile/config.txt /etc
```

Packager is using a Docker feature called Named contexts. Simply it adds
`--build-context package-context=<package_context_root>` option to `docker build` command. The
`package-context` is the name of the context (used with `--from`). The Named context is supported
by multiple Dockerfile commands. More information can be find in
[Docker documentation](https://docs.docker.com/build/concepts/context/).
