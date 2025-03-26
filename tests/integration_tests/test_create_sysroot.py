import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager, is_tracked, does_app_support_image, prepare_packages


def test_01_crete_sysroot(test_image, packager_binary, context, test_repo, test_sysroot):
    """TODO"""
    package = "test_package"
    app = "io-module"

    prepare_packages([package], context)
    # print(does_app_support_image(app, test_image))

    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=package,
        expected_result=True,
    )

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=app,
        expected_result=True,
    )
    assert is_tracked(app, test_repo, "app")

    run_packager(
        packager_binary,
        "create-sysroot",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=app,
        expected_result=True,
    )
    # assert is_tracked(app, test_repo, "app")


def test_02_create_sysroot_inconsistent_image_names(test_image, packager_binary, context, test_repo, test_sysroot):
    """TODO"""
    return
