import pytest
import subprocess
import os
import test_utils.test_utils as test_utils


# @pytest.fixture(scope="session")
# def package_to_image_placer_binary():
#     """Compile the Go application binary."""
#     package_to_image_placer_binary = "../package-to-image-placer"
#     source_file = "../ImageToPackagePlacer.go"

#     subprocess.run(["go", "build", "-o", package_to_image_placer_binary, source_file], check=True)

#     yield package_to_image_placer_binary  # Pass the binary path to tests


@pytest.fixture(autouse=True)
def clean_up_between_tests():
    """Set up the environment before each test is run and clean up after each test is run."""
    test_data_dir = "test_data"
    # remove any previous test data
    if os.path.exists(test_data_dir):
        test_utils.remove_dir(test_data_dir)

    os.makedirs(test_data_dir, exist_ok=True)

    yield

    subprocess.run(["sudo", "losetup", "-D"], check=True)
    # Clean up
    # subprocess.run(["sudo", "rm", "-rf", test_data_dir], check=True)


@pytest.fixture(scope="session", autouse=True)
def setup_environment():
    """Set up the environment before any tests are run."""
    # try:
         # subprocess.run(["go", "build", "-o", "../package-to-image-placer", "../ImageToPackagePlacer.go"], check=True)

    # except FileNotFoundError:
    #     pytest.fail("Please make sure the app is built before running the tests.")

    print("Checking if all system utilities are installed...")
    try:
        # subprocess.run(["go", "version"], check=True, capture_output=True, text=True)
        # subprocess.run(["touch", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["rm", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["mkfs", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["mkdir", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["zip", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["unzip", "-v"], check=True, capture_output=True, text=True)
        # subprocess.run(["mount", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["umount", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["sync", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["fallocate", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["parted", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["losetup", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["mkfs.vfat", "--help"], check=True, capture_output=True, text=True)
        # subprocess.run(["mkfs.ext4", "-V"], check=True, capture_output=True, text=True)
        # subprocess.run(["diff", "--version"], check=True, capture_output=True, text=True)
        # subprocess.run(["libguestfs-test-tool"], check=True, capture_output=True, text=True)

    except FileNotFoundError:
        pytest.fail("Please make sure all required system utilities are installed.")

    yield
