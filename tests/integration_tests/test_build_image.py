import pytest
import threading

from test_utils.test_utils import run_packager, does_image_exist, check_stdout, test_config
from test_utils.common import PackagerExpectedResult


def test_01_show_help(packager_binary):
    """Show image help"""

    for mode in ["build-image", "build-package", "build-app", "create-sysroot"]:
        result = run_packager(packager_binary, mode, help=True)
        stdout, stderr = result.communicate()

        assert "usage: bringauto packager" in stdout.lower()
        assert "error" not in stdout.lower()
        assert stderr == ""


def test_02_build_image(test_image, packager_binary, context):
    """Test build image"""

    run_packager(
        packager_binary,
        "build-image",
        context=context,
        image_name=test_image,
        expected_result=PackagerExpectedResult.SUCCESS,
    )
    assert does_image_exist(test_image)


def test_03_build_all_images(packager_binary, context, request):
    """Test build all images. Only runs if --image is set to 'all'"""
    if request.config.getoption("--image") != "all":
        pytest.skip("Skipping test because --image is not set to 'all'")

    run_packager(
        packager_binary, "build-image", context=context, all=True, expected_result=PackagerExpectedResult.SUCCESS
    )

    for image in test_config["test_images"]:
        assert does_image_exist(image)


def test_04_build_multiple_images_at_once(packager_binary, context):
    """Test building two images at once"""

    def build_image(image_name):
        run_packager(
            packager_binary,
            "build-image",
            context=context,
            image_name=image_name,
            expected_result=PackagerExpectedResult.SUCCESS,
        )

    threads = []
    images = ["fedora40", "fedora41"]
    for image in images:
        t = threading.Thread(target=build_image, args=(image,))
        t.start()
        threads.append(t)

    for t in threads:
        t.join()

    assert does_image_exist("fedora40")
    assert does_image_exist("fedora41")
