from enum import Enum


class PackagerReturnCode(Enum):
    """
    Enum to represent the return codes from the packager.
    """

    SUCCESS = 0
    DEFAULT_ERROR = 1  # Not specified errors
    CMD_LINE_ERROR = 2  # Cmd line parsing errors
    CONTEXT_ERROR = 3  # Context consistency errors
    GIT_LFS_ERROR = 4  # Git Lfs and Context comparison consistency errors
    BUILD_ERROR = 5  # Build image, App or Package errors
    PACKAGE_MISSING_DEPENDENCY_ERROR = 6  # Package dependency is not on sysroot error
    CREATING_SYSROOT_ERROR = 7  # Creating sysroot errors
    OVERWRITE_FILE_IN_SYSROOT_ERROR = 8  # Overwriting files in sysroot error


class PackagerExpectedResult(Enum):
    """
    Enum to represent the expected results from the packager.
    """

    SUCCESS = True
    FAILURE = False
    NOT_APPLICABLE = None  # Used when the test is not applicable for the given context
    SKIPPED = "SKIPPED"  # Used when the test is skipped due to some condition
    CREATING_SYSROOT = "CREATING_SYSROOT"  # Used when the test is creating a sysroot
