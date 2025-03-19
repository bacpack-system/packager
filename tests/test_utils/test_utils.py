import subprocess
import os
import git
import json
import docker
import shutil
from time import sleep


test_config = {
    "test_data_examples": "test_data/example",
    "test_apps": os.path.abspath("test_data/example/app"),
    "test_dockers": os.path.abspath("test_data/example/docker"),
    "test_packages": os.path.abspath("test_data/example/package"),
    "packager_binary": os.path.abspath("../bap-builder/bap-builder"),
    "test_repo": os.path.abspath("test_data/test_repo"),
    "install_sysroot": os.path.abspath("install_sysroot"),
    "test_packages_source": os.path.abspath("test_data/test_packages"),
}
test_config["test_images"] = sorted(os.listdir(test_config["test_dockers"]))


def init_test_repo():
    if os.path.exists(test_config["test_repo"]):
        shutil.rmtree(test_config["test_repo"])

    os.makedirs(test_config["test_repo"])
    git.Repo.init(test_config["test_repo"])
    return test_config["test_repo"]


def clean():
    if os.path.exists(test_config["install_sysroot"]):
        shutil.rmtree(test_config["install_sysroot"])

    if os.path.exists(test_config["test_packages"]):
        shutil.rmtree(test_config["test_packages"])
    os.makedirs(test_config["test_packages"])
    init_test_repo()

    pass


def does_image_exist(image: str) -> bool:
    """Check if the image exists in the local Docker registry."""
    client = docker.from_env()
    try:
        client.images.get(image)
        return True
    except docker.errors.ImageNotFound:
        return False


def check_stdout(stdout: str, expected_result: bool):
    stdout = stdout.lower()
    if expected_result:
        assert "build ok" in stdout
        assert "error" not in stdout
        assert "failed to" not in stdout
    else:
        assert "error" in stdout
        assert "failed to build" in stdout


def prepare_packages(packages: list[str]):
    for package in packages:
        package_path = os.path.join(test_config["test_packages"], package)
        if os.path.exists(package_path):
            shutil.rmtree(package_path)
        shutil.copytree(os.path.join(test_config["test_packages_source"], package), package_path)


def is_package_tracked(package_name: str, repo_path: str) -> bool:
    """Check if the package is tracked in the repository."""
    repo = git.Repo(repo_path)
    try:
        files_in_last_commit = repo.git.log("--diff-filter=A", "--name-only", "--pretty=format:").splitlines()
    except git.exc.GitCommandError:
        files_in_last_commit = []

    for path in files_in_last_commit:
        if package_name in path.split("/"):
            if package_name in path.split("/")[-1]:
                return True
            else:
                return False

    return False


def run_packager(
    packager_binary: str,
    mode: str,
    context: str = None,
    image_name: str = None,
    output_dir: str = None,
    name: str = None,
    build_deps: bool = False,
    build_deps_on: bool = False,
    build_deps_on_recursive: bool = False,
    help: bool = False,
    all: bool = False,
    expected_result: bool = None,
) -> subprocess.CompletedProcess:
    """TODO"""

    parameters = [packager_binary, mode]

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

    if output_dir:
        parameters.append("--output-dir")
        parameters.append(output_dir)

    if name:
        parameters.append("--name")
        parameters.append(name)

    if build_deps:
        parameters.append("--build-deps")

    if build_deps_on:
        parameters.append("--build-deps-on")

    if build_deps_on_recursive:
        parameters.append("--build-deps-on-recursive")

    print(parameters)

    result = subprocess.Popen(
        parameters,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    stdout, stdin = result.communicate()

    print(stdout)
    print(stdin)

    if expected_result is not None:
        check_stdout(stdout, expected_result)

    assert result.returncode == 0
    return result
