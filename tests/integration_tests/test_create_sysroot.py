import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, is_tracked, does_app_support_image, prepare_packages


def test_01_crete_sysroot(test_image, packager_binary, context, test_repo, test_sysroot):
    """TODO"""
    package = "test_package_1"
    app = "io-module"

    prepare_packages([package])
    # print(does_app_support_image(app, test_image))

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
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=app,
        expected_result=True,
    )
    assert is_tracked(app, test_repo, "app")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name=test_image,
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=True,
    )
    print("done")
    # assert is_tracked(app, test_repo, "app")


def test_02_create_sysroot_inconsistent_image_names(packager_binary, context, test_repo, test_sysroot):
    """TODO"""
    package = "test_package_1"
    app = "io-module"

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
        "build-app",
        context=context,
        image_name="fedora40",
        output_dir=test_repo,
        name=app,
        expected_result=True,
    )
    assert is_tracked(app, test_repo, "app")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora41",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=True,
    )
    print("done")


def test_03_create_sysroot_from_empty_repo(packager_binary, context, test_repo, test_sysroot):
    """TODO"""
    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora41",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=False,
    )
    print("done")


def test_04_create_sysroot_from_repo_with_packages_for_different_images(
    packager_binary, context, test_repo, test_sysroot
):
    """TODO"""
    package = "test_package_1"

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
    print("done")


def test_05_create_sysroot_from_all_packages(packager_binary, context, test_repo, test_sysroot):
    """TODO"""
    package = "test_package_1"

    prepare_packages([package])
    # print(does_app_support_image(app, test_image))

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
