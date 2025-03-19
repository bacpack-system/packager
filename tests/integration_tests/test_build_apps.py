import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, check_stdout, is_package_tracked, prepare_packages


def test_01_build_app(test_image, packager_binary, context, test_repo):
    """TODO"""
    app = "io-module"

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=app,
    )


def test_02_build_all_apps(test_image, packager_binary, context, test_repo):
    """TODO"""
    apps = ["io-module", "io-module2"]

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        all=True,
    )
