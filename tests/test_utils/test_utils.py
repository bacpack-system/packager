from logging import exception
import subprocess
import os
import git
import json
import docker
import shutil
from time import sleep


test_config = {
    "test_data_examples": os.path.abspath("test_data/example"),
    "test_apps": os.path.abspath("test_data/example/app"),
    "test_dockers": os.path.abspath("../example/docker"),
    "test_packages": os.path.abspath("test_data/example/package"),
    "packager_binary": os.path.abspath("../bap-builder/bap-builder"),
    "test_repo": os.path.abspath("test_data/test_repo"),
    "test_sysroot": os.path.abspath("test_data/test_sysroot"),
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


def init_test_sysroot():
    if os.path.exists(test_config["test_sysroot"]):
        shutil.rmtree(test_config["test_sysroot"])

    os.makedirs(test_config["test_sysroot"])
    return test_config["test_sysroot"]


def clean_sysroot():
    if os.path.exists(test_config["install_sysroot"]):
        shutil.rmtree(test_config["install_sysroot"])


def clean():
    clean_sysroot()
    if os.path.exists(test_config["test_packages"]):
        shutil.rmtree(test_config["test_packages"])
    os.makedirs(test_config["test_packages"])
    init_test_repo()


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
        assert "build ok" in stdout or "creating sysroot directory from packages" in stdout
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


def is_tracked(name: str, repo_path: str, type: str) -> bool:
    """Check if the package is tracked in the repository."""
    repo = git.Repo(repo_path)
    try:
        files_in_last_commit = repo.git.log("--diff-filter=A", "--name-only", "--pretty=format:").splitlines()
    except git.exc.GitCommandError:
        files_in_last_commit = []

    if type == "app":
        source_path = os.path.join(test_config["test_apps"], name)
        debug_suffix = "d_"
        release_suffix = "_"
    elif type == "package":
        source_path = os.path.join(test_config["test_packages"], name)
        debug_suffix = "d-dev_"
        release_suffix = "-dev_"
    else:
        raise ValueError("Invalid type")

    files = os.listdir(source_path)
    packages_to_detect = len(files)

    for file in files:
        if "debug" in file:
            test_name = name + debug_suffix
        else:
            test_name = name + release_suffix

        for path in files_in_last_commit:
            if name in path.split("/"):
                if test_name in path.split("/")[-1]:
                    packages_to_detect -= 1

    return packages_to_detect == 0


def get_nested(data, keys, default=None):
    for key in keys:
        if isinstance(data, dict) and key in data:
            data = data[key]
        else:
            return default
    return data


def does_app_support_image(app: str, image: str) -> bool:
    """Check if the app supports the given image."""
    try:
        for apps in os.listdir(os.path.join(test_config["test_apps"], app)):
            with open(os.path.join(test_config["test_apps"], app, apps)) as file:
                metadata = json.load(file)
                if image not in get_nested(metadata, ["DockerMatrix", "ImageNames"], []):
                    return False
        return True

    except FileNotFoundError:
        raise FileNotFoundError(f"App {app} does not exist.")
    except Exception as e:
        print(e)


def run_packager(
    packager_binary: str,
    mode: str,
    context: str = None,
    image_name: str = None,
    output_dir: str = None,
    name: str = None,
    sysroot_dir: str = None,
    git_lfs: str = None,
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

    if sysroot_dir:
        parameters.append("--sysroot-dir")
        parameters.append(sysroot_dir)

    if git_lfs:
        parameters.append("--git-lfs")
        parameters.append(git_lfs)

    if name:
        parameters.append("--name")
        parameters.append(name)

    if build_deps:
        parameters.append("--build-deps")

    if build_deps_on:
        parameters.append("--build-deps-on")

    if build_deps_on_recursive:
        parameters.append("--build-deps-on-recursive")

    print("\033[95m\nRunning command:", " ".join(parameters), "\033[0m")

    result = subprocess.Popen(
        parameters,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    stdout, stderr = result.communicate()

    print("-" * 10, "Stdout", "-" * 10)
    print(stdout)
    print("-" * 10, "Stderr", "-" * 10)
    print(stderr)

    if expected_result is not None:
        check_stdout(stdout, expected_result)

    assert result.returncode == 0
    return result
