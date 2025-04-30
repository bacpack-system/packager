import subprocess
import os
from time import sleep

from test_utils.test_utils import (
    run_packager,
    is_tracked,
    does_app_support_image,
    clean_sysroot,
)


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
        expected_result=None if not does_app_support_image(app, test_image) else True,
    )
    if does_app_support_image(app, test_image):
        assert is_tracked(app, test_repo, "app")
    else:
        assert not is_tracked(app, test_repo, "app")


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
        expected_result=None if not does_app_support_image(apps[0], test_image) else True,
    )

    clean_sysroot()

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=apps[1],
        expected_result=None if not does_app_support_image(apps[1], test_image) else True,
    )
    for app in apps:
        if does_app_support_image(app, test_image):
            assert is_tracked(app, test_repo, "app")
        else:
            assert not is_tracked(app, test_repo, "app")


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
        expected_result=None,
    )
    for app in apps:
        if does_app_support_image(app, test_image):
            assert is_tracked(app, test_repo, "app")
        else:
            assert not is_tracked(app, test_repo, "app")
