import pytest

from test_utils.test_utils import run_packager, does_image_exist, check_stdout, test_config


def test_01_show_help(packager_binary):
    """TODO"""

    for mode in ["build-image", "build-package", "build-app", "create-sysroot"]:
        result = run_packager(packager_binary, mode, help=True)
        stdout, stderr = result.communicate()

        assert "usage: bringauto packager" in stdout.lower()
        assert "error" not in stdout.lower()
        assert stderr == ""


def test_02_build_image(test_image, packager_binary, context):
    """TODO"""

    run_packager(packager_binary, "build-image", context=context, image_name=test_image, expected_result=True)
    assert does_image_exist(test_image)


def test_03_build_all_images(packager_binary, context, request):
    """TODO"""
    if request.config.getoption("--image") != "all":
        pytest.skip("Skipping test because --image is not set to 'all'")

    run_packager(packager_binary, "build-image", context=context, all=True, expected_result=True)

    for image in test_config["test_images"]:
        assert does_image_exist(image)
