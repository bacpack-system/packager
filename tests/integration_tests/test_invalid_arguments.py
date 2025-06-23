import subprocess
import pathlib
import os

from test_utils.test_utils import run_packager

from test_utils.common import PackagerReturnCode


def test_01_run_without_command(packager_binary):
    """TODO"""
    result = run_packager(packager_binary, "", expected_returncode=PackagerReturnCode.CMD_LINE_ERROR)
    stdout = result.communicate()[0]

    assert "ERROR" in stdout


def test_02_run_with_invalid_command(packager_binary):
    """TODO"""
    result = run_packager(packager_binary, "invalid-command", expected_returncode=PackagerReturnCode.CMD_LINE_ERROR)
    stdout = result.communicate()[0]

    assert "ERROR" in stdout


def test_03_run_with_nonexisting_image(packager_binary):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    result = run_packager(
        packager_binary,
        "build-image",
        context=context,
        image_name="nonexisting_image",
        expected_returncode=PackagerReturnCode.DEFAULT_ERROR,
    )
    stdout = result.communicate()[0]

    assert "ERROR" in stdout
    assert "Failed to build Docker image:" in stdout


def test_04_run_with_nonexisting_package(test_image, packager_binary, test_repo):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    result = run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        name="nonexisting_package",
        output_dir=test_repo,
        expected_returncode=PackagerReturnCode.DEFAULT_ERROR,
    )
    stdout = result.communicate()[0]

    assert "ERROR" in stdout


def test_05_run_with_nonexisting_app(test_image, packager_binary, test_repo):
    """TODO"""
    context = os.path.abspath(os.path.join("test_data", "example"))

    result = run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        name="nonexisting_app",
        output_dir=test_repo,
        expected_returncode=PackagerReturnCode.DEFAULT_ERROR,
    )
    stdout = result.communicate()[0]

    assert "ERROR" in stdout
