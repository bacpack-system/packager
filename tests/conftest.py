import pytest
import subprocess
import os
import shutil
import docker

packager_binary_path = "../bap-builder/bap-builder"
test_data_examples = "test_data/example"
test_apps = os.path.join(test_data_examples, "app")
test_dockers = os.path.join(test_data_examples, "docker")
test_packages = os.path.join(test_data_examples, "package")


def get_available_images():
    return sorted(os.listdir(test_dockers))


@pytest.fixture(scope="session", autouse=True)
def setup_environment(request):
    """Set up the environment before any tests are run."""
    if request.config.getoption("--remove_images"):
        # TODO test this feature
        client = docker.from_env()
        for image in get_available_images():
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


@pytest.fixture
def test_images(request):
    all_images = get_available_images()
    image_param = request.config.getoption("--image")
    if image_param == "all":
        return all_images
    elif image_param in all_images:
        return [image_param]
    else:
        raise ValueError(f"Invalid image: {image_param}. Available images: {all_images}")


@pytest.fixture(scope="session")
def packager_binary():
    """Compile the Go application binary."""

    subprocess.run(["go", "get", "bringauto/bap-builder"], check=True)
    subprocess.run(["go", "build", "-o", packager_binary_path, "../bap-builder"], check=True)

    yield packager_binary_path  # Pass the binary path to tests


@pytest.fixture(autouse=True)
def clean_up_between_tests():
    """Set up the environment before each test is run and clean up after each test is run."""

    yield
