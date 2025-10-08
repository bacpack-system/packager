
# Docker Container Requirements

Each image that we can use for build our dependencies

## SSH Server

- SSH server must be enabled on standard port (22)
- `permitRootLogin` must be enabled in the `sshd` configuration
- password for user `root` must be `1234`

## CMake

Required if any Package built for given image uses CMake build system.

- CMake >= 3.21 must be installed in the system and reachable for user `root`

## Meson

Required if any Package built for given image uses Meson build system.

- Meson must be installed in the system and reachable for user `root`
- Ninja must be installed in the system and reachable for user `root`

## Bash

- Standard `bash` utility must be installed and reachable for user `root`

## os-release and uname

`/etc/os-release` file and `uname` are used to construct platform string.

`/etc/os-release` must include `ID` and `VERSION_ID` fields.

`uname` must support

- `-m` - machine. For example "x86_64"

# Host system

Docker container forward port 22 of the sshd daemon in the container to the
specified port (1122 by default) of the host system.
