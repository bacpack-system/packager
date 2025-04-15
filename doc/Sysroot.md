# Sysroot

BAP is creating its own sysroot directories when building Packages. The separated sysroot
directories are created for both debug and release Packages. All Package build files are copied to
these directories. The desired behaviour is to ensure that no Package build files are being
ovewritten by another Package build files. To ensure this, following mechanisms have been
introduced to BAP.

## Sysroot consistency mechanisms

- At the start of `build-package` and `build-app` command the both debug and release sysroots are
checked. If it isn't empty, the warning is printed.

- During Package/App builds, the Package/App build files are copied to the sysroot directory
(`install_sysroot`). If any of the file is already present in the sysroot, the error is printed
that the Package tries to overwrite files in sysroot, which would corrupt consistency of sysroot
directory. If Package doesn't try to overwrite any files, the build proceeds and Package files are
added to the Package Repository.

- Copied Package informations are added to built Packages file (more below). When the
`build-package` command with `--build-deps-on` option is used, it is expected that the Package with
its dependencies are already in sysroot. If it is not, the error is printed and build fails.

- When `create-sysroot` command is used, all Packages in Package Repository for given target platform
files are copied to new sysroot directory. Because of the sysroot consistency mechanism this new
sysroot will also be consistent.

- When `build-app` command with `--all` option is used, the sysroot directory is deleted after each
build (with Apps the sysroot does not have to be shared between builds). If single App is being
build, the sysroot is not deleted so the user can check the sysroot after build.

## Built Packages

The `built_packages.json` file in `install_sysroot` directory contains already built Packages in
sysroot. When building a Package, it is first checked if the Package is already built in sysroot.
If it is, the build is skipped, else it continues. It is also used when building with
`--build-deps-on` flag for check that needed Packages are already built.

Built Package in `install_sysroot` is identified with these values:

- Package name also with `d` suffix and `lib` prefix
- Directory name in sysroot - this is for resolution of Package builds for specific image
- Git URL of the Package
- Git commit hash which is obtained after the Package is git cloned to docker container

The git informations ensures that if the Package git repository is changed (commit hash or git URL
changes) the Package will be build again. Note that the build will probably fail, because the
Package would probably overwrite its own files in sysroot.

If the user wants to force build of Package already built in sysroot, the `install_sysroot`
directory must be deleted.

## Notes

- The `install_sysroot` directory is not being deleted when building Packages (for a backup
reason), it is only deleted when building all Apps
