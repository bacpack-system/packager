
# Platform String

Platform string is a string in form `MACHINE` - `DISTRO_NAME` - `DISTRO_VERSION`. It is used for
Package zip filenames and paths in Package Repository and sysroot directory.

## Construction

The platform string is constructed from three parts:

- `MACHINE` - machine architecture. For example `x86-64` or `arm64`
- `DISTRO_NAME` - name of the distribution. For example `debian` or `fedora`
- `DISTRO_VERSION` - version of the distribution. For example `12` or `41`

For `MACHINE` the `uname -m` output is used. The output is converted to lowercase and the `_` is
replaced with `-`. The uname is executed in the container.

For `DISTRO_NAME` and `DISTRO_VERSION` the `/etc/os-release` file is used. The file is copied from
the container and parsed.

## Examples

- `x86-64-debian-12`
- `arm64-fedora-41`
