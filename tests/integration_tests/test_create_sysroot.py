import subprocess
import os
from time import sleep
import pytest

from test_utils.test_utils import (
    run_packager,
    is_tracked,
    does_app_support_image,
    does_package_support_image,
    prepare_packages,
    check_if_package_is_in_sysroot,
)

from test_utils.common import PackagerReturnCode


def test_01_create_sysroot(test_image, packager_binary, context, test_repo, test_sysroot):
    """Build package and app, then create sysroot"""
    packages = ["test_package_1", "test_package_2"]
    app = "io-module"

    prepare_packages(packages)

    if not does_app_support_image(app, test_image) or not all(
        does_package_support_image(pkg, test_image) for pkg in packages
    ):
        pytest.skip(f"Skipping test because {app} or {packages} does not support {test_image}")

    # run_packager(
    # packager_binary,
    # "build-app",
    # context=context,
    # image_name=test_image,
    # output_dir=test_repo,
    # name=app,
    # expected_result=True,
    # )
    # assert is_tracked(app, test_repo, "app")
    #
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[1],
        build_deps=True,
        expected_result=True,
    )
    assert is_tracked(packages[1], test_repo, "package")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name=test_image,
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=True,
    )


def test_02_create_sysroot_with_package_on_two_different_images(packager_binary, context, test_repo, test_sysroot):
    """Test creating sysroot with a package built on two different images"""
    package = "test_package_1"
    # app = "io-module"

    prepare_packages([package])
    # print(does_app_support_image(app, test_image))

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora41",
        output_dir=test_repo,
        name=package,
        expected_result=True,
    )
    assert is_tracked(package, test_repo, "package", os_path="fedora/41")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora40",
        output_dir=test_repo,
        name=package,
        expected_result=True,
    )
    assert is_tracked(package, test_repo, "package", os_path="fedora/40")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora40",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=True,
    )


def test_03_create_sysroot_from_empty_repo(packager_binary, context, test_repo, test_sysroot):
    """Create sysroot from empty repo"""
    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora41",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=False,
        expected_returncode=PackagerReturnCode.CREATING_SYSROOT_ERROR,
    )


def test_04_create_sysroot_from_repo_with_packages_for_different_images(
    packager_binary, context, test_repo, test_sysroot
):
    """Create sysroot for image which is not built in the repo"""
    package = "test_package_1"

    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora41",
        output_dir=test_repo,
        name=package,
        expected_result=True,
    )
    assert is_tracked(package, test_repo, "package")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora40",
        output_dir=test_repo,
        name=package,
        expected_result=True,
    )
    assert is_tracked(package, test_repo, "package")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="ubuntu2204",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=False,
    )


def test_05_create_sysroot_from_all_packages(packager_binary, context, test_repo, test_sysroot):
    """Build all packages and create sysroot from all packages"""
    package = "test_package_1"

    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora41",
        output_dir=test_repo,
        all=True,
        expected_result=True,
    )

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora41",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=True,
    )


def test_06_check_data_in_sysroot(test_image, packager_binary, context, test_repo, test_sysroot):
    """Check"""
    package = "test_package_2"
    app = "io-module"

    prepare_packages([package])
    if not does_app_support_image(app, test_image) or not does_package_support_image(package, test_image):
        pytest.skip(f"Skipping test because {app} or {package} does not support {test_image}")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=True,
    )
    assert is_tracked(package, test_repo, "package")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name=test_image,
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=True,
    )
    files = [
        "release/bin/curl",
        "release/bin/curl-config",
        "release/include/curl/curl.h",
        "release/include/curl/curlver.h",
        "release/include/curl/easy.h",
        "release/include/curl/mprintf.h",
        "release/include/curl/multi.h",
        "release/include/curl/options.h",
        "release/include/curl/stdcheaders.h",
        "release/include/curl/system.h",
        "release/include/curl/system.h",
        "release/include/curl/typecheck-gcc.h",
        "release/include/curl/urlapi.h",
        "release/bin/curl",
        "release/bin/curl-config",
        "release/lib64/libcurl.so",
        "release/lib64/pkgconfig/libcurl.pc",
    ]
    assert check_if_package_is_in_sysroot(test_sysroot, files)
