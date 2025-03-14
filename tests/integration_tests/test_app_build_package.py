import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, get_available_images, check_stdout, is_package_tracked


def test_app_build_package(test_images, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_1"

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_images[0],
        output_dir=test_repo,
        package_name=package,
    )
    stdout = result.communicate()[0]

    assert result.returncode == 0
    check_stdout(stdout, True)


def test_app_build_package_with_dependency(test_images, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_images[0],
        output_dir=test_repo,
        package_name=depends_on_package,
    )
    stdout = result.communicate()[0]

    assert result.returncode == 0
    check_stdout(stdout, True)
    assert not is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_images[0],
        output_dir=test_repo,
        package_name=package,
    )
    stdout = result.communicate()[0]

    assert result.returncode == 0
    check_stdout(stdout, True)
    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)


def test_app_build_package_dependency_with_build_deps(test_images, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_images[0],
        output_dir=test_repo,
        package_name=package,
        build_deps=True,
    )
    stdout = result.communicate()[0]
    assert result.returncode == 0
    check_stdout(stdout, True)

    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)


def test_app_build_package_dependency_with_build_deps_on(test_images, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_images[0],
        output_dir=test_repo,
        package_name=depends_on_package,
    )
    stdout = result.communicate()[0]
    assert result.returncode == 0
    check_stdout(stdout, True)

    assert not is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_images[0],
        output_dir=test_repo,
        package_name=depends_on_package,
        build_deps_on=True,
    )
    stdout = result.communicate()[0]
    assert result.returncode == 0
    check_stdout(stdout, True)

    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)
