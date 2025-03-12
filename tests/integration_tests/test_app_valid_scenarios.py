import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager
import pytest


def test_app_show_help(test_images, packager_binary):
    """Tests if the package_to_image_placer will fail when the paths are invalid"""

    for mode in ["build-image", "build-package", "build-app", "create-sysroot"]:
        result = run_packager(packager_binary, mode, help=True)
        stdout, stderr = result.communicate()

        assert result.returncode == 0
        assert "usage: bringauto packager" in stdout.lower()
        assert "error" not in stdout.lower()
        assert stderr == ""


def test_app_build_image(test_images, packager_binary):
    """Tests if the package_to_image_placer will fail when the paths are invalid"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    for image in test_images:
        result = run_packager(packager_binary, "build-image", context=context, image=image)
        stdout = result.communicate()[0]

        assert result.returncode == 0
        assert "Build OK" in stdout


def test_app_build_all_images(test_images, packager_binary):
    """Tests if the package_to_image_placer will fail when the paths are invalid"""
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
    assert "Can't build image" not in stdout
