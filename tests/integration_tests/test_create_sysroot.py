import pytest

from test_utils.test_utils import (
    run_packager,
    is_tracked,
    does_package_support_image,
    does_packages_support_image,
    prepare_packages,
    check_if_package_is_in_sysroot,
)

from test_utils.common import PackagerReturnCode, PackagerExpectedResult


def test_01_create_sysroot(test_image, packager_binary, context, test_repo, test_sysroot):
    """Build package and then create sysroot"""
    packages = ["test-package-a", "test-package-b"]

    prepare_packages(packages)

    if not all(does_package_support_image(pkg, test_image) for pkg in packages):
        pytest.skip(f"Skipping test because {packages} does not support {test_image}")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[1],
        build_deps=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(packages[1], test_repo, "package")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name=test_image,
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=PackagerExpectedResult.CREATING_SYSROOT,
    )


def test_02_create_sysroot_with_package_on_two_different_images(
    packager_binary, context, test_repo, test_sysroot
):
    """Test creating sysroot with a package built on two different images"""
    package = "test-package-a"

    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="ubuntu2404",
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(package, test_repo, "package", os_path="ubuntu/24.04")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora43",
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(package, test_repo, "package", os_path="fedora/43")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora43",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=PackagerExpectedResult.CREATING_SYSROOT,
    )


def test_03_create_sysroot_from_empty_repo(packager_binary, context, test_repo, test_sysroot):
    """Create sysroot from empty repo"""
    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora43",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=False,
        expected_returncode=PackagerReturnCode.CREATING_SYSROOT_ERROR,
    )


def test_04_create_sysroot_from_repo_with_packages_for_different_images(
    packager_binary, context, test_repo, test_sysroot
):
    """Create sysroot for image which is not built in the repo"""
    package = "test-package-a"

    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="debian13",
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(package, test_repo, "package", os_path="debian/13")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora43",
        output_dir=test_repo,
        name=package,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(package, test_repo, "package", os_path="fedora/43")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="ubuntu2404",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=False,
        expected_returncode=PackagerReturnCode.CREATING_SYSROOT_ERROR,
    )


def test_05_create_sysroot_from_all_packages(packager_binary, context, test_repo, test_sysroot):
    """Build all packages and create sysroot from all packages"""
    packages = ["test-package-a", "test-package-b", "test-package-c", "test-package-d", "test-package-e"]

    prepare_packages(packages)

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name="fedora43",
        output_dir=test_repo,
        all=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name="fedora43",
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=PackagerExpectedResult.CREATING_SYSROOT,
    )


def test_06_check_data_in_sysroot(test_image, packager_binary, context, test_repo, test_sysroot):
    """Check that expected files from the package are present in the created sysroot"""
    packages = ["test-package-a", "test-package-b"]

    prepare_packages(packages)

    if not does_packages_support_image(packages, test_image):
        pytest.skip(f"Skipping test because {packages} does not support {test_image}")

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=packages[1],
        build_deps=True,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert is_tracked(packages[1], test_repo, "package")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name=test_image,
        sysroot_dir=test_sysroot,
        git_lfs=test_repo,
        expected_result=PackagerExpectedResult.CREATING_SYSROOT,
    )
    files = [
        "release/lib/libtest-package-a-shared.so",
        "release/lib/libtest-package-b-shared.so",
        "release/include/pack_a.hpp",
        "release/include/pack_b.hpp",
    ]

    assert check_if_package_is_in_sysroot(test_sysroot, files)
