import pytest
import subprocess
import docker
from test_utils.test_utils import test_config, init_test_repo, clean, init_test_sysroot, run_packager
from test_utils.common import PackagerExpectedResult


def pytest_sessionstart(session):
    """Run before any tests, shows as setup in output."""
    print("=" * 50, "SETUP PHASE", "=" * 50)

    check_or_remove_images(session)

    print("SETUP: Building Packager binary")

    try:
        subprocess.run(["go", "version"], check=True, capture_output=True, text=True)

    except FileNotFoundError:
        pytest.fail("Please make sure all required system utilities are installed.")
    
    subprocess.run(
        ["go", "build", "-o", test_config["packager_binary"], "../cmd/bap-builder"],
        check=True,
    )
    print("DONE")
    print("SETUP: Building Docker images")
    
    # Build all images
    run_packager(
        test_config["packager_binary"],
        "build-image",
        context=test_config["test_context"],
        all=True,
        expected_result=PackagerExpectedResult.SUCCESS,
        show_output=False,
    )

    print("DONE")


def pytest_sessionfinish(session):
    clean()


def pytest_addoption(parser):
    parser.addoption(
        "--image",
        action="store",
        default="fedora43",
        help="The image to use for testing",
    )
    parser.addoption(
        "--remove_images",
        action="store_true",
        default=False,
        help="Remove build docker images before running tests so that they are rebuilt",
    )


def check_or_remove_images(session):
    client = docker.from_env()
    if session.config.getoption("--remove_images"):
        # TODO test this feature
        print("SETUP: Removing old Docker images")
        
        for image in test_config["test_images"]:
            try:
                client.images.remove(image, force=True)
            except docker.errors.ImageNotFound:
                pass
        print("DONE")
    else:
        existing_images = []
        for image in test_config["test_images"]:
            try:
                client.images.get(image)
            except docker.errors.ImageNotFound:
                pass
            else:
                existing_images.append(image)
        if len(existing_images) > 0:
            pytest.exit(f"Images {existing_images} already exists. Use --remove_images to remove it.")


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
    return test_config["test_context"]


@pytest.fixture
def test_repo():
    return init_test_repo()


@pytest.fixture
def test_sysroot():
    return init_test_sysroot()


@pytest.fixture(scope="session")
def packager_binary():
    """Compile the Go application binary."""

    return test_config["packager_binary"]  # Pass the binary path to tests


@pytest.fixture(autouse=True)
def clean_up_between_tests():
    """Set up the environment before each test is run and clean up after each test is run."""

    clean()

    yield

    # if os.path.exists(test_repo_path):
    #     shutil.rmtree(test_repo_path)
