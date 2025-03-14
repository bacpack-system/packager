import subprocess
import pathlib
import os

from test_utils.test_utils import run_packager


def test_run_without_command(test_images, packager_binary):
    """TODO"""
    result = run_packager(packager_binary, "")
    stdout = result.communicate()[0]

    assert result.returncode == 0
    assert "ERROR" in stdout


def test_run_with_invalid_command(test_images, packager_binary):
    """TODO"""
    result = run_packager(packager_binary, "invalid-command")
    stdout = result.communicate()[0]

    assert result.returncode == 0
    assert "ERROR" in stdout


def test_run_with_nonexisting_image(test_images, packager_binary):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    result = run_packager(packager_binary, "build-image", context=context, image_name="nonexisting_image")
    stdout = result.communicate()[0]

    assert result.returncode == 0
    assert "ERROR" in stdout
    assert "Failed to build Docker image:" in stdout


def test_run_with_nonexisting_package(test_images, packager_binary):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    result = run_packager(
        packager_binary, "build-package", context=context, image_name=test_images[0], package_name="nonexisting_package"
    )
    stdout = result.communicate()[0]

    assert result.returncode == 0
    assert "ERROR" in stdout
