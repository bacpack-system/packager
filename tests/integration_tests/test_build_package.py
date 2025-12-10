from test_utils.test_utils import (
    run_packager,
    is_tracked,
    prepare_packages,
    does_package_support_image,
    does_packages_support_image,
)

from test_utils.common import PackagerReturnCode, PackagerExpectedResult


def test_01_build_package(test_image, packager_binary, context, test_repo):
    """Test building a package with no dependencies."""
    package = "test-package-a"
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
    package = "test-package-b"
    depends_on_package = "test-package-a"
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


def test_03_build_package_dependency_with_build_deps(
    test_image, packager_binary, context, test_repo
):
    """Test build package with deps on."""
    package = "test-package-b"
    depends_on_package = "test-package-a"
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


def test_04_build_multiple_package_dependency_with_build_deps(
    test_image, packager_binary, context, test_repo
):
    """Test build package with multiple dependencies using deps on."""
    packages = ["test-package-a", "test-package-b", "test-package-c", "test-package-d"]
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


def test_05_build_package_dependency_with_build_deps_on(
    test_image, packager_binary, context, test_repo
):
    """First build dependency and then build package with deps on"""
    package = "test-package-b"
    depends_on_package = "test-package-a"
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


def test_06_build_multiple_package_dependency_with_build_deps_on(
    test_image, packager_binary, context, test_repo
):
    """Test build package with multiple dependencies using deps on."""
    packages = [
        "test-package-a",
        "test-package-b",
        "test-package-c",
        "test-package-d",
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
    packages = ["test-package-a", "test-package-b", "test-package-c", "test-package-d"]
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
    package = "test-package-circ-dep-a"
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
    package = "test-package-circ-dep-a"
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


def test_10_has_itself_as_dependency_build_deps_on_recursive(
    test_image, packager_binary, context, test_repo
):
    """Test build package with itself as dependency using build_deps_on_recursive."""
    package = "test-package-circ-dep-a"
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
    packages = ["test-package-circ-dep-b", "test-package-circ-dep-c", "test-package-circ-dep-d"]
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
    packages = ["test-package-circ-dep-b", "test-package-circ-dep-c", "test-package-circ-dep-d"]
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


def test_13_circular_dependencies_deps_on_recursive(
    test_image, packager_binary, context, test_repo
):
    """Test build package with circular dependencies using deps on recursive."""
    packages = ["test-package-circ-dep-b", "test-package-circ-dep-c", "test-package-circ-dep-d"]
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
        assert not is_tracked(
            package, test_repo, "package"
        ), f"Package {package} should not be tracked but is."


def test_14_fork_dependencies_deps(test_image, packager_binary, context, test_repo):
    """Test build package with fork dependencies using deps."""
    packages = [
        "test-package-e",
        "test-package-a",
        "test-package-b",
        "test-package-c",
        "test-package-d",
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
        build_deps=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    for package in packages:
        assert is_tracked(
            package, test_repo, "package"
        ), f"Package {package} should be tracked but is not."


def test_15_fork_dependencies_deps_on(test_image, packager_binary, context, test_repo):
    """Test build package with fork dependencies using deps on."""
    packages = [
        "test-package-a",
        "test-package-b",
        "test-package-c",
        "test-package-d",
        "test-package-e",
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

    assert is_tracked(
        packages[0], test_repo, "package"
    ), f"Package {packages[0]} should be tracked but is not."

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
        assert is_tracked(
            package, test_repo, "package"
        ), f"Package {package} should be tracked but is not."


def test_16_fork_dependencies_deps_on_recursive(test_image, packager_binary, context, test_repo):
    """Test build package with fork dependencies using deps on recursive."""
    packages = [
        "test-package-a",
        "test-package-b",
        "test-package-c",
        "test-package-d",
        "test-package-e",
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
    package = "test-package-a-17"
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
    assert is_tracked(
        package, test_repo, "package"
    ), f"Package {package} should be tracked but is not."


def test_18_only_release(test_image, packager_binary, context, test_repo):
    """Test building a package with only release version."""
    package = "test-package-a-18"
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
    assert is_tracked(
        package, test_repo, "package"
    ), f"Package {package} should be tracked but is not."


def test_19_missing_release_debug_packages_build_deps(
    test_image, packager_binary, context, test_repo
):
    """Test building a package where the dependency is missing debug version"""
    packages = ["test-package-a-17", "test-package-b-17"]
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
        assert not is_tracked(
            package, test_repo, "package"
        ), f"Package {package} should not be tracked but is."


def test_20_missing_release_debug_packages_build_deps_on(
    test_image, packager_binary, context, test_repo
):
    """Test building a package where the deps on is missing release version"""
    packages = ["test-package-a-17", "test-package-b-17"]
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
        assert not is_tracked(
            package, test_repo, "package"
        ), f"Package {package} should not be tracked but is."


def test_21_build_packages_with_no_images(test_image, packager_binary, context, test_repo):
    """Build a package with no supported images. It should cause a context error."""
    package = "test-package-a-21"
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


def test_22_build_packages_where_dependency_is_not_supported(
    test_image, packager_binary, context, test_repo
):
    """Dependent package is not supported by the image, but the package itself is supported, so it should not be built."""
    package = "test-package-b"
    dependents_on_package = "test-package-a-22"
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


def test_23_build_packages_where_package_is_not_supported(
    test_image, packager_binary, context, test_repo
):
    """Dependent package is supported but the package itself is not supported, so nothing should be built."""
    package = "test-package-b-23"
    dependents_on_package = "test-package-a"
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
    package = "test-package-a"
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


def test_25_invalid_revision(test_image, packager_binary, context, test_repo):
    """Test building a package with an invalid revision."""
    package = "test-package-a-25"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.BUILD_ERROR,
    )

def test_26_meson_build(test_image, packager_binary, context, test_repo):
    """Test building a package using Meson build system."""
    package = "test-package-f"
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

    assert is_tracked(package, test_repo, "package")

def test_27_cmake_invalid_define(test_image, packager_binary, context, test_repo):
    """Test building a package with an CMake define with invalid characters."""
    package = "test-package-a-27"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

def test_28_meson_invalid_define(test_image, packager_binary, context, test_repo):
    """Test building a package with an Meson define with invalid characters."""
    package = "test-package-f-28"
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
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

def test_29_meson_invalid_option(test_image, packager_binary, context, test_repo):
    """Test building a package with an Meson option with invalid characters."""
    package = "test-package-f-29"
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
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

def test_30_meson_empty_define(test_image, packager_binary, context, test_repo):
    """Test building a package with an Meson define with empty name."""
    package = "test-package-f-30"
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
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

def test_31_meson_empty_option(test_image, packager_binary, context, test_repo):
    """Test building a package with an Meson option with empty name."""
    package = "test-package-f-31"
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
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

def test_32_cmake_empty_define(test_image, packager_binary, context, test_repo):
    """Test building a package with an CMake define with empty name."""
    package = "test-package-a-32"
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
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )

def test_33_cmake_wrong_define_type(test_image, packager_binary, context, test_repo):
    """Test building a package with an CMake define with non-string value."""
    package = "test-package-a-33"
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
        expected_result=PackagerExpectedResult.FAILURE,
        expected_returncode=PackagerReturnCode.CONTEXT_ERROR,
    )
