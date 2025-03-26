import subprocess
import os
from time import sleep

from test_utils.test_utils import (
    run_packager,
    does_image_exist,
    check_stdout,
    is_tracked,
    prepare_packages,
    does_app_support_image,
)


def test_01_build_app(test_image, packager_binary, context, test_repo):
    """TODO"""
    app = "io-module"
    # app = "mission-module"
    print(does_app_support_image(app, test_image))

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


def test_02_build_multiple_apps(test_image, packager_binary, context, test_repo):
    """TODO"""
    apps = ["io-module", "mission-module"]
    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=apps[0],
        expected_result=True,
    )
    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=apps[1],
        expected_result=True,
    )
    for app in apps:
        assert is_tracked(app, test_repo, "app")


def test_03_build_all_apps(test_image, packager_binary, context, test_repo):
    """TODO"""
    apps = ["io-module", "mission-module"]
    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        all=True,
    )
    for app in apps:
        assert is_tracked(app, test_repo, "app")
