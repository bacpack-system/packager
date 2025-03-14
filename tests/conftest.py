import pytest
import subprocess
import os
import shutil
import docker
from test_utils.test_utils import test_config, init_test_repo, clean
from git import Repo


@pytest.fixture(scope="session", autouse=True)
def setup_environment(request):
    """Set up the environment before any tests are run."""
    if request.config.getoption("--remove_images"):
        # TODO test this feature
        client = docker.from_env()
        for image in test_config["test_images"]:
            try:
                client.images.remove(image, force=True)
            except docker.errors.ImageNotFound:
                pass

    try:
        subprocess.run(["go", "version"], check=True, capture_output=True, text=True)

    except FileNotFoundError:
        pytest.fail("Please make sure all required system utilities are installed.")

    yield


def pytest_addoption(parser):
    parser.addoption("--image", action="store", default="ubuntu2204", help="The image to use for testing")
    parser.addoption(
        "--remove_images",
        action="store_true",
        default=False,
        help="Remove build docker images before running tests so that they are rebuilt",
    )


@pytest.fixture(params=test_config["test_images"])
def test_image(request):
    image_param = request.config.getoption("--image")
    if image_param == "all":
        return request.param
    elif image_param == request.param:
        return request.param
    else:
        pytest.skip(f"Skipping test for image {request.param}")


@pytest.fixture
def context():
    return os.path.abspath(os.path.join("test_data", "example"))


@pytest.fixture
def test_repo():
    return init_test_repo()


@pytest.fixture(scope="session")
def packager_binary():
    """Compile the Go application binary."""

    subprocess.run(["go", "get", "bringauto/bap-builder"], check=True)
    subprocess.run(["go", "build", "-o", test_config["packager_binary"], "../bap-builder"], check=True)

    yield test_config["packager_binary"]  # Pass the binary path to tests


@pytest.fixture(autouse=True)
def clean_up_between_tests():
    """Set up the environment before each test is run and clean up after each test is run."""

    clean()

    yield

    # if os.path.exists(test_repo_path):
    #     shutil.rmtree(test_repo_path)
