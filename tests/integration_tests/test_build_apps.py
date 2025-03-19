import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, does_image_exist, check_stdout, is_package_tracked, prepare_packages


def test_01_build_app(test_image, packager_binary, context, test_repo):
    """TODO"""
    package = "test_package_1"
    prepare_packages([package])

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
    )
