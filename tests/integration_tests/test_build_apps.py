import threading
from time import sleep
import socket
import pytest

from test_utils.test_utils import (
    run_packager,
    is_tracked,
    does_app_support_image,
    clean_sysroot,
)

from test_utils.common import PackagerExpectedResult, PackagerReturnCode


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
        expected_result=(
            PackagerExpectedResult.SUCCESS
            if does_app_support_image(app, test_image)
            else PackagerExpectedResult.FAILURE
        ),
        expected_returncode=(
            PackagerReturnCode.SUCCESS if does_app_support_image(app, test_image) else PackagerReturnCode.DEFAULT_ERROR
        ),
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
        expected_result=(
            PackagerExpectedResult.SUCCESS
            if does_app_support_image(apps[0], test_image)
            else PackagerExpectedResult.NOT_APPLICABLE
        ),
        expected_returncode=(
            PackagerReturnCode.SUCCESS
            if does_app_support_image(apps[0], test_image)
            else PackagerReturnCode.DEFAULT_ERROR
        ),
    )

    clean_sysroot()

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=apps[1],
        expected_result=(
            PackagerExpectedResult.SUCCESS
            if does_app_support_image(apps[1], test_image)
            else PackagerExpectedResult.NOT_APPLICABLE
        ),
        expected_returncode=(
            PackagerReturnCode.SUCCESS
            if does_app_support_image(apps[1], test_image)
            else PackagerReturnCode.DEFAULT_ERROR
        ),
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
        expected_result=(
            PackagerExpectedResult.SUCCESS
            if all(does_app_support_image(app, test_image) for app in apps)
            else PackagerExpectedResult.NOT_APPLICABLE
        ),
        expected_returncode=(
            PackagerReturnCode.SUCCESS
            if all(does_app_support_image(app, test_image) for app in apps)
            else PackagerReturnCode.DEFAULT_ERROR
        ),
    )
    if all(does_app_support_image(app, test_image) for app in apps):
        for app in apps:
            assert is_tracked(app, test_repo, "app")
    else:
        for app in apps:
            assert not is_tracked(app, test_repo, "app")


def test_04_build_all_apps_when_port_1122_is_used(test_image, packager_binary, context, test_repo):
    """Test building all apps when port 1122 is already used"""
    apps = ["io-module", "mission-module"]
    if not all(does_app_support_image(app, test_image) for app in apps):
        pytest.skip(f"Apps does not support image {test_image}")

    # Seize port 1122 before running the test
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        sock.bind(("localhost", 1122))
        sock.listen(1)

        run_packager(
            packager_binary,
            "build-app",
            context=context,
            image_name=test_image,
            output_dir=test_repo,
            all=True,
            expected_result=PackagerExpectedResult.SUCCESS,
        )
        for app in apps:
            assert is_tracked(app, test_repo, "app")
    finally:
        sock.close()


def test_05_build_app_in_parallel(test_image, packager_binary, context, test_repo):
    """Test building an app in parallel threads"""

    apps = ["io-module", "mission-module"]
    if not all(does_app_support_image(app, test_image) for app in apps):
        pytest.skip(f"Apps does not support image {test_image}")

    def build_app(app):
        run_packager(
            packager_binary,
            "build-app",
            context=context,
            image_name=test_image,
            output_dir=test_repo,
            name=app,
            expected_result=(PackagerExpectedResult.SUCCESS),
        )

    threads = []
    for app in apps:
        thread = threading.Thread(target=build_app, args=(app,))
        threads.append(thread)
        thread.start()

    for thread in threads:
        thread.join()

    for app in apps:
        assert is_tracked(app, test_repo, "app")
