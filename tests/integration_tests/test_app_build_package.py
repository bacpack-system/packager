import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, check_stdout, is_package_tracked, prepare_packages


def test_01_build_package(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_1"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=package,
        expected_result=True,
    )


def test_02_build_package_with_dependency(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"
    prepare_packages([package, depends_on_package])

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
        package_name=package,
        expected_result=True,
    )
    assert is_package_tracked(package, test_repo)
    assert is_package_tracked(depends_on_package, test_repo)


def test_03_build_package_dependency_with_build_deps(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"
    prepare_packages([package, depends_on_package])

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


def test_04_build_multiple_package_dependency_with_build_deps(test_image, packager_binary, context, test_repo):
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    prepare_packages(packages)

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


def test_05_build_package_dependency_with_build_deps_on(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_2"
    depends_on_package = "test_package_1"
    prepare_packages([package, depends_on_package])

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


def test_06_build_multiple_package_dependency_with_build_deps_on(test_image, packager_binary, context, test_repo):
    packages = ["test_package_1_06", "test_package_2_06", "test_package_3_06", "test_package_4_06"]
    prepare_packages(packages)

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


def test_07_build_multiple_package_dependency_with_build_deps_on_recursive(
    test_image, packager_binary, context, test_repo
):
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    prepare_packages(packages)

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


def test_08_has_itself_as_dependency_build_deps(test_image, packager_binary, context, test_repo):
    package = "test_package_5"
    prepare_packages([package])

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


def test_09_has_itself_as_dependency_build_deps_on(test_image, packager_binary, context, test_repo):
    package = "test_package_5"
    prepare_packages([package])

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


def test_10_has_itself_as_dependency_build_deps_on_recursive(test_image, packager_binary, context, test_repo):
    package = "test_package_5"
    prepare_packages([package])

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


def test_11_circular_dependencies_deps(test_image, packager_binary, context, test_repo):
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    prepare_packages(packages)

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


def test_12_circular_dependencies_deps_on(test_image, packager_binary, context, test_repo):
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    prepare_packages(packages)

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

    assert not is_package_tracked(packages[0], test_repo)
    assert not is_package_tracked(packages[1], test_repo)
    assert is_package_tracked(packages[2], test_repo)


def test_13_circular_dependencies_deps_on_recursive(test_image, packager_binary, context, test_repo):
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    prepare_packages(packages)

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
        build_deps_on_recursive=True,
        expected_result=False,
    )

    assert not is_package_tracked(packages[0], test_repo)
    assert not is_package_tracked(packages[1], test_repo)
    assert is_package_tracked(packages[2], test_repo)


def test_14_fork_dependencies_deps(test_image, packager_binary, context, test_repo):
    packages = ["test_package_9", "test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    prepare_packages(packages)

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[0],
        build_deps=True,
        expected_result=True,
    )

    for package in packages:
        assert is_package_tracked(package, test_repo)


def test_15_fork_dependencies_deps_on(test_image, packager_binary, context, test_repo):
    """FIXME"""

    packages = [
        "test_package_1",
        "test_package_2",
        "test_package_3",
        "test_package_4",
        "test_package_9",
        "test_package_5",
        "test_package_6",
        "test_package_7",
        "test_package_8",
    ]
    prepare_packages(packages)

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

    for package in packages[:5]:
        assert is_package_tracked(package, test_repo)

    for package in packages[5:]:
        assert not is_package_tracked(package, test_repo)


def test_16_fork_dependencies_deps_on_recursive(test_image, packager_binary, context, test_repo):
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4", "test_package_9"]
    prepare_packages(packages)

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


def test_17_only_debug(test_image, packager_binary, context, test_repo):
    package = "test_package_1_17"
    prepare_packages([package])

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


def test_18_only_release(test_image, packager_binary, context, test_repo):
    package = "test_package_2_17"
    prepare_packages([package])

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


def test_19_missing_release_debug_packages_build_deps(test_image, packager_binary, context, test_repo):
    packages = ["test_package_1_17", "test_package_2_17"]
    prepare_packages(packages)

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[1],
        build_deps=True,
        expected_result=False,
    )

    for package in packages:
        assert not is_package_tracked(package, test_repo)


def test_20_missing_release_debug_packages_build_deps_on(test_image, packager_binary, context, test_repo):
    """FIXME"""
    packages = ["test_package_1_17", "test_package_2_17"]
    prepare_packages(packages)

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        package_name=packages[1],
        build_deps_on=True,
        expected_result=False,
    )

    assert not is_package_tracked(packages[0], test_repo)
    assert not is_package_tracked(packages[1], test_repo)
