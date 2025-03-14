import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, check_stdout
import pytest


def test_app_show_help(packager_binary):
    """TODO"""

    for mode in ["build-image", "build-package", "build-app", "create-sysroot"]:
        result = run_packager(packager_binary, mode, help=True)
        stdout, stderr = result.communicate()

        assert "usage: bringauto packager" in stdout.lower()
        assert "error" not in stdout.lower()
        assert stderr == ""


def test_app_build_image(test_image, packager_binary, context):
    """TODO"""

    run_packager(packager_binary, "build-image", context=context, image_name=test_image, expected_result=True)
    assert does_image_exist(test_image)
