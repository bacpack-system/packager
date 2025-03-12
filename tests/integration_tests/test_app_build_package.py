import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, get_available_images, check_stdout


def test_app_build_package(test_images, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_3"

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

    # for package in os.listdir(os.path.join("test_data", "example", "package")):
    #     assert os.path.exists(os.path.join("test_data", "example", "package", package))


def test_app_build_package_with_dependency(test_images, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_1"
    depends_on_package = "spdlog"

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
