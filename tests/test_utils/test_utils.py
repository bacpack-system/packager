import subprocess
import os
import shutil
import json
import docker
from time import sleep

test_data_examples = "test_data/example"
test_apps = os.path.join(test_data_examples, "app")
test_dockers = os.path.join(test_data_examples, "docker")
test_packages = os.path.join(test_data_examples, "package")


def does_image_exist(image: str) -> bool:
    """Check if the image exists in the local Docker registry."""
    client = docker.from_env()
    try:
        client.images.get(image)
        return True
    except docker.errors.ImageNotFound:
        return False


def get_available_images():
    return sorted(os.listdir(test_dockers))


def check_stdout(stdout: str, expected_result: bool):
    stdout = stdout.lower()
    if expected_result:
        assert "build ok" in stdout
        assert "error" not in stdout
        assert "failed to" not in stdout
    else:
        assert "error" in stdout.upper()
        assert "failed to build docker image:" in stdout


def run_packager(
    package_binary: str,
    mode: str,
    context: str = None,
    image_name: str = None,
    help: bool = False,
    all: bool = False,
) -> subprocess.CompletedProcess:
    """TODO"""

    parameters = [package_binary, mode]

    if help:
        parameters.append("--help")

    if all:
        parameters.append("--all")

    if context:
        parameters.append("--context")
        parameters.append(context)

    if image_name:
        parameters.append("--image-name")
        parameters.append(image_name)

    # print(parameters)

    print(parameters)

    result = subprocess.Popen(
        parameters,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    stdout, stdin = result.communicate()

    # this outputs can be inspected when running pytest with -s flag
    print(stdout)
    print(stdin)

    return result
