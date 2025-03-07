import subprocess
import os
from time import sleep

from test_utils.test_utils import run_packager


def test_invalid_file_paths(packager_binary):
    """Tests if the package_to_image_placer will fail when the paths are invalid"""

    rc = run_packager(package_binary=packager_binary)

    print(rc.returncode)

    # print(rc.stdout)
    # print(rc.stderr)
