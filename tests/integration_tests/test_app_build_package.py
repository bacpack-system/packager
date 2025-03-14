import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, check_stdout, is_package_tracked


def test_app_build_package(test_image, packager_binary, context, test_repo, expected_result=True):
    """TODO"""
    package = "test_package_1"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
    )


def test_app_build_package_with_dependency(test_image, packager_binary, context, test_repo, expected_result=True):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=depends_on_package,
    )
    assert not is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
        expected_result=True,
    )
    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)


def test_app_build_package_dependency_with_build_deps(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
        build_deps=True,
        expected_result=True,
    )
    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)


def test_app_build_multiple_package_dependency_with_build_deps(test_image, packager_binary, context, test_repo):
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[-1],
        build_deps=True,
        expected_result=True,
    )
    for package in packages:
        assert is_package_tracked(package, test_repo)


def test_app_build_package_dependency_with_build_deps_on(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=depends_on_package,
        expected_result=True,
    )
    assert not is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=depends_on_package,
        build_deps_on=True,
        expected_result=True,
    )
    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)


def test_app_build_multiple_package_dependency_with_build_deps_on(test_image, packager_binary, context, test_repo):
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        expected_result=True,
    )

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        build_deps_on=True,
        expected_result=True,
    )

    assert is_package_tracked(packages[0], test_repo)
    assert is_package_tracked(packages[1], test_repo)
    assert not is_package_tracked(packages[2], test_repo)
    assert not is_package_tracked(packages[3], test_repo)


def test_app_build_multiple_package_dependency_with_build_deps_on_recursive(
    test_image, packager_binary, context, test_repo
):
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        expected_result=True,
    )
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        build_deps_on_recursive=True,
        expected_result=True,
    )
    for package in packages:
        assert is_package_tracked(package, test_repo)


def test_app_has_itself_as_dependency_build_deps(test_image, packager_binary, context, test_repo):
    package = "test_package_5"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
        build_deps=True,
        expected_result=False,
    )
    assert not is_package_tracked(package, test_repo)


def test_app_has_itself_as_dependency_build_deps_on(test_image, packager_binary, context, test_repo):
    package = "test_package_5"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
        build_deps_on=True,
        expected_result=False,
    )
    assert not is_package_tracked(package, test_repo)


def test_app_has_itself_as_dependency_build_deps_on_recursive(test_image, packager_binary, context, test_repo):
    package = "test_package_5"

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
        build_deps_on_recursive=True,
        expected_result=False,
    )
    assert not is_package_tracked(package, test_repo)


def test_app_circular_dependencies_deps(test_image, packager_binary, context, test_repo):
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        build_deps=True,
        expected_result=False,
    )

    for package in packages:
        assert not is_package_tracked(package, test_repo)


def test_app_circular_dependencies_deps_on(test_image, packager_binary, context, test_repo):
    packages = ["test_package_6", "test_package_7", "test_package_8"]

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[-1],
        expected_result=True,
    )

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        build_deps_on=True,
        expected_result=False,
    )

    for package in packages:
        assert not is_package_tracked(package, test_repo)
