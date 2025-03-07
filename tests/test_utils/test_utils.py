import subprocess
import os
import shutil
import json
from time import sleep


def run_packager(
    package_binary: str,
) -> subprocess.CompletedProcess:
    """TODO"""

    parameters = [package_binary]

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
