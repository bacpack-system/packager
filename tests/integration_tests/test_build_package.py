import subprocess
import os
from time import sleep

from test_utils.test_utils import (
    run_packager,
    does_image_exist,
    check_stdout,
    is_tracked,
    prepare_packages,
    does_package_support_image,
    does_packages_support_image,
)

from test_utils.common import PackagerReturnCode, PackagerExpectedResult


def test_01_build_package(test_image, packager_binary, context, test_repo):
    """Test building a package with no dependencies."""
    package = "test_package_1"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=(
            PackagerExpectedResult.NOT_APPLICABLE
            if not does_package_support_image(package, test_image)
            else PackagerExpectedResult.SUCCESS
        ),
    )


def test_02_build_package_with_dependency(test_image, packager_binary, context, test_repo):
    """Test building a package with a dependency. Fist build the dependency and then the package."""
    package = "test_package_2"
    depends_on_package = "test_package_1"
    prepare_packages([package, depends_on_package])

    if not does_package_support_image(depends_on_package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=depends_on_package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert not is_tracked(package, test_repo, "package")
    assert is_tracked(depends_on_package, test_repo, "package")

    if not does_package_support_image(package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    assert is_tracked(package, test_repo, "package")
    assert is_tracked(depends_on_package, test_repo, "package")


def test_03_build_package_dependency_with_build_deps(test_image, packager_binary, context, test_repo):
    """Test build package with deps on."""
    package = "test_package_2"
    depends_on_package = "test_package_1"
    prepare_packages([package, depends_on_package])

    if not does_packages_support_image([package, depends_on_package], test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        build_deps=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    assert is_tracked(package, test_repo, "package")
    assert is_tracked(depends_on_package, test_repo, "package")


def test_04_build_multiple_package_dependency_with_build_deps(test_image, packager_binary, context, test_repo):
    """Test build package with multiple dependencies using deps on."""
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[-1],
        build_deps=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    for package in packages:
        assert is_tracked(package, test_repo, "package")


def test_05_build_package_dependency_with_build_deps_on(test_image, packager_binary, context, test_repo):
    """First build dependency and then build package with deps on"""
    package = "test_package_2"
    depends_on_package = "test_package_1"
    prepare_packages([package, depends_on_package])

    if not does_packages_support_image([package, depends_on_package], test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=depends_on_package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert not is_tracked(package, test_repo, "package")
    assert is_tracked(depends_on_package, test_repo, "package")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=depends_on_package,
        build_deps_on=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(package, test_repo, "package")
    assert is_tracked(depends_on_package, test_repo, "package")


def test_06_build_multiple_package_dependency_with_build_deps_on(test_image, packager_binary, context, test_repo):
    """Test build package with multiple dependencies using deps on."""
    packages = ["test_package_1_06", "test_package_2_06", "test_package_3_06", "test_package_4_06"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps_on=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    assert is_tracked(packages[0], test_repo, "package")
    assert is_tracked(packages[1], test_repo, "package")
    assert not is_tracked(packages[2], test_repo, "package")
    assert not is_tracked(packages[3], test_repo, "package")


def test_07_build_multiple_package_dependency_with_build_deps_on_recursive(
    test_image, packager_binary, context, test_repo
):
    """Test build package with multiple dependencies using deps on recursive."""
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps_on_recursive=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    for package in packages:
        assert is_tracked(package, test_repo, "package")


def test_08_has_itself_as_dependency_build_deps(test_image, packager_binary, context, test_repo):
    """Test build package with itself as dependency using deps."""
    package = "test_package_5"
    prepare_packages([package])

    if not does_package_support_image(package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        build_deps=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )
    assert not is_tracked(package, test_repo, "package")


def test_09_has_itself_as_dependency_build_deps_on(test_image, packager_binary, context, test_repo):
    """Test build package with itself as dependency using build_deps_on."""
    package = "test_package_5"
    prepare_packages([package])

    if not does_package_support_image(package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        build_deps_on=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )
    assert not is_tracked(package, test_repo, "package")


def test_10_has_itself_as_dependency_build_deps_on_recursive(test_image, packager_binary, context, test_repo):
    """Test build package with itself as dependency using build_deps_on_recursive."""
    package = "test_package_5"
    prepare_packages([package])

    if not does_package_support_image(package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        build_deps_on_recursive=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )
    assert not is_tracked(package, test_repo, "package")


def test_11_circular_dependencies_deps(test_image, packager_binary, context, test_repo):
    """Test build package with circular dependencies using deps."""
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    for package in packages:
        assert not is_tracked(package, test_repo, "package")


def test_12_circular_dependencies_deps_on(test_image, packager_binary, context, test_repo):
    """Test build package with circular dependencies using deps on."""
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[-1],
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps_on=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    assert all(not is_tracked(package, test_repo, "package") for package in packages)


def test_13_circular_dependencies_deps_on_recursive(test_image, packager_binary, context, test_repo):
    """Test build package with circular dependencies using deps on recursive."""
    packages = ["test_package_6", "test_package_7", "test_package_8"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[-1],
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps_on_recursive=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    for package in packages:
        assert not is_tracked(package, test_repo, "package"), f"Package {package} should not be tracked but is."


def test_14_fork_dependencies_deps(test_image, packager_binary, context, test_repo):
    """Test build package with fork dependencies using deps."""
    packages = ["test_package_9", "test_package_1", "test_package_2", "test_package_3", "test_package_4"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    for package in packages:
        assert is_tracked(package, test_repo, "package"), f"Package {package} should be tracked but is not."


def test_15_fork_dependencies_deps_on(test_image, packager_binary, context, test_repo):
    """Test build package with fork dependencies using deps on."""
    packages = [
        "test_package_1",
        "test_package_2",
        "test_package_3",
        "test_package_4",
        "test_package_9",
    ]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    assert is_tracked(packages[0], test_repo, "package"), f"Package {packages[0]} should be tracked but is not."

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps_on=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    for package in packages:
        assert is_tracked(package, test_repo, "package"), f"Package {package} should be tracked but is not."


def test_16_fork_dependencies_deps_on_recursive(test_image, packager_binary, context, test_repo):
    """Test build package with fork dependencies using deps on recursive."""
    packages = ["test_package_1", "test_package_2", "test_package_3", "test_package_4", "test_package_9"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    assert is_tracked(packages[0], test_repo, "package")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[0],
        build_deps_on_recursive=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    for package in packages:
        assert is_tracked(package, test_repo, "package")


def test_17_only_debug(test_image, packager_binary, context, test_repo):
    """Test building a package with only debug version."""
    package = "test_package_1_17"
    prepare_packages([package])

    if not does_package_support_image(package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(package, test_repo, "package"), f"Package {package} should be tracked but is not."


def test_18_only_release(test_image, packager_binary, context, test_repo):
    """Test building a package with only release version."""
    packages = ["test_package_1_17", "test_package_2_17"]
    prepare_packages(packages)

    if not does_package_support_image(packages[-1], test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[-1],
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )
    for package in packages:
        assert not is_tracked(package, test_repo, "package"), f"Package {package} should not be tracked but is."


def test_19_missing_release_debug_packages_build_deps(test_image, packager_binary, context, test_repo):
    """Test building a package where the dependency is missing debug version"""
    packages = ["test_package_1_17", "test_package_2_17"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[1],
        build_deps=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    for package in packages:
        assert not is_tracked(package, test_repo, "package"), f"Package {package} should not be tracked but is."


def test_20_missing_release_debug_packages_build_deps_on(test_image, packager_binary, context, test_repo):
    """Test building a package where the deps on is missing release version"""
    packages = ["test_package_1_17", "test_package_2_17"]
    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[1],
        build_deps_on=True,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    for package in packages:
        assert not is_tracked(package, test_repo, "package"), f"Package {package} should not be tracked but is."


def test_21_build_packages_with_no_images(test_image, packager_binary, context, test_repo):
    """Build a package with no supported images. It should cause a context error."""
    package = "test_package_1_21"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        image_name=test_image,
        context=context,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.NOT_APPLICABLE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    assert not is_tracked(package, test_repo, "package")


def test_22_build_packages_where_dependency_is_not_supported(test_image, packager_binary, context, test_repo):
    """Dependent package is not supported by the image, but the package itself is supported, so it should not be built."""
    package = "test_package_1_22"
    dependents_on_package = "test_package_2_22"
    prepare_packages([package, dependents_on_package])

    if not does_package_support_image(package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.FAILURE,
        build_deps=True,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    assert not is_tracked(package, test_repo, "package")
    assert not is_tracked(dependents_on_package, test_repo, "package")


def test_23_build_packages_where_package_is_not_supported(test_image, packager_binary, context, test_repo):
    """Dependent package is supported but the package itself is not supported, so nothing should be built."""
    package = "test_package_2_23"
    dependents_on_package = "test_package_1_23"
    prepare_packages([package, dependents_on_package])

    if not does_package_support_image(dependents_on_package, test_image):
        return

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.FAILURE,
        build_deps=True,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

    assert not is_tracked(package, test_repo, "package")
    assert not is_tracked(dependents_on_package, test_repo, "package")


def test_24_build_same_package_twice(test_image, packager_binary, context, test_repo):
    """Build the same package twice in a row to check if the second build is skipped."""
    package = "test_package_1"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=(
            PackagerExpectedResult.SUCCESS
            if does_package_support_image(package, test_image)
            else PackagerExpectedResult.NOT_APPLICABLE
        ),
    )

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=(
            PackagerExpectedResult.NOT_APPLICABLE
            if not does_package_support_image(package, test_image)
            else PackagerExpectedResult.SUCCESS
        ),
    )
