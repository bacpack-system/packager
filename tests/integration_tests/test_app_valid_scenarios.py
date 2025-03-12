import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, get_available_images, check_stdout
import pytest


def test_app_show_help(test_images, packager_binary):
    """TODO"""

    for mode in ["build-image", "build-package", "build-app", "create-sysroot"]:
        result = run_packager(packager_binary, mode, help=True)
        stdout, stderr = result.communicate()

        assert result.returncode == 0
        assert "usage: bringauto packager" in stdout.lower()
        assert "error" not in stdout.lower()
        assert stderr == ""


def test_app_build_image(test_images, packager_binary):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    for image in test_images:
        result = run_packager(packager_binary, "build-image", context=context, image_name=image)
        stdout = result.communicate()[0]

        assert result.returncode == 0
        check_stdout(stdout, True)
        assert does_image_exist(image)


def test_app_build_all_images(test_images, packager_binary):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    if len(test_images) <= 1:
        # Run it only when all packages are specified
        pytest.skip(
            "Skipping test_app_build_all_images because not all images were selected. Use --image=all to specify all images."
        )
        return

    result = run_packager(packager_binary, "build-image", context=context, all=True)
    stdout = result.communicate()[0]

    assert result.returncode == 0
    check_stdout(stdout, True)

    for image in get_available_images():
        assert does_image_exist(image)
