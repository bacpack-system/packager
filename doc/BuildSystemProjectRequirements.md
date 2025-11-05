
# Build System Project Requirements

At this moment Packager supports these build systems for C/C++ projects:

- CMake
- Meson

## Requirements

- Project must be able to be installed by `make install` or `meson install`
- Project must NOT override `CMAKE_INSTALL_PREFIX` CMake variable or `prefix` Meson option - it's used for the project installation to a given directory and Package creation. If you override it the build fail!
- Project must NOT override `CMAKE_PREFIX_PATH` CMake variable or `cmake-prefix-path` Meson option - it's used for finding dependencies in the build sysroot. If you override it the build fail!
