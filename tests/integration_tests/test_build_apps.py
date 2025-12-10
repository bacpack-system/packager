import socket
import pytest

from test_utils.test_utils import (
    run_packager,
    is_tracked,
    does_app_support_image,
    clean_sysroot,
    prepare_packages,
)

from test_utils.common import PackagerExpectedResult, PackagerReturnCode


def test_01_build_app(test_image, packager_binary, context, test_repo):
    """Test building a single app"""
    app = "test-app-a"
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
            PackagerReturnCode.SUCCESS
            if does_app_support_image(app, test_image)
            else PackagerReturnCode.DEFAULT_ERROR
        ),
    )
    if does_app_support_image(app, test_image):
        assert is_tracked(app, test_repo, "app")
    else:
        assert not is_tracked(app, test_repo, "app")


def test_02_build_multiple_apps(test_image, packager_binary, context, test_repo):
    """Test building multiple apps at once"""
    apps = ["test-app-a", "test-app-b"]
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


def build_app_3_deps(test_image, packager_binary, context, test_repo):
    deps = ["test-package-e", "test-package-a", "test-package-b", "test-package-c", "test-package-d"]

    prepare_packages(deps)

    # Build required dependency Packages
    run_packager(
        packager_binary,
        "build-package",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        build_deps=True,
        name=deps[0],
        expected_result=PackagerExpectedResult.SUCCESS,
        expected_returncode=PackagerReturnCode.SUCCESS,
    )

def test_03_build_app_with_dependencies(test_image, packager_binary, context, test_repo):
    """Test building app which has dependencies on Packages in Package Repository"""
    app = "test-app-c"
    build_app_3_deps(test_image, packager_binary, context, test_repo)

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        name=app,
        use_local_repo=True,
        expected_result=(
            PackagerExpectedResult.SUCCESS
            if does_app_support_image(app, test_image)
            else PackagerExpectedResult.NOT_APPLICABLE
        ),
        expected_returncode=(
            PackagerReturnCode.SUCCESS
            if does_app_support_image(app, test_image)
            else PackagerReturnCode.DEFAULT_ERROR
        ),
    )


def test_04_build_all_apps(test_image, packager_binary, context, test_repo):
    """Test building all apps"""
    apps = ["test-app-a", "test-app-b", "test-app-c"]
    build_app_3_deps(test_image, packager_binary, context, test_repo)

    run_packager(
        packager_binary,
        "build-app",
        context=context,
        image_name=test_image,
        output_dir=test_repo,
        all=True,
        use_local_repo=True,
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


def test_05_build_all_apps_when_port_1122_is_used(
    test_image, packager_binary, context, test_repo
):
    """Test building all apps when port 1122 is already used"""
    apps = ["test-app-a", "test-app-b", "test-app-c"]
    if not all(does_app_support_image(app, test_image) for app in apps):
        pytest.skip(f"Apps does not support image {test_image}")
    
    build_app_3_deps(test_image, packager_binary, context, test_repo)

    # Seize port 1122 before running the test
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    try:
        sock.bind(("localhost", 1122))
        sock.listen(1)

        run_packager(
            packager_binary,
            "build-app",
            context=context,
            image_name=test_image,
            output_dir=test_repo,
            port=1123,
            all=True,
            use_local_repo=True,
            expected_result=PackagerExpectedResult.SUCCESS,
        )
        for app in apps:
            assert is_tracked(app, test_repo, "app")
    except OSError:
        pytest.skip("Port 1122 is already in use, skipping test")
    finally:
        sock.close()
