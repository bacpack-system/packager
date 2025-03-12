import subprocess
import os
import shutil
import json
from time import sleep


def run_packager(
    package_binary: str,
    mode: str,
    context: str = None,
    image: str = None,
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

    if image:
        parameters.append("--name")
        parameters.append(image)

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
