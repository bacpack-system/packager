# Tests documentation

## TODO

- [ ] Add test documentation
- [ ] Test if docker remove works
- [ ] remove dependency on port 1122
- [ ] check repo consistency on the start of process

## Requirements

The tests share the same requirements as the program itself and additionally depend on standard Linux utilities and Python 3. Written in Python, they leverage pytest for execution. It is recommended to use a virtual environment for managing Python dependencies.

```bash
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Before running the tests, ensure that the application can be built successfully. Follow the build instructions provided in the root `README.md` file.

<!-- ---
**⚠️ WARNING**  
To ensure the tests run correctly, please take note of the following requirements:  

1. **Root Privileges:** The application requires root privileges during testing. Ensure you have the **root password** available when prompted.
2. **Non-standard Test Behavior:** If tests are interrupted or exhibit non-standard behavior, you may need to manually unmount any devices used during the tests. In some cases, a system restart might be necessary to restore normal operation.

--- -->

## Running tests

To run tests use the following command:

```bash
pytest
```

To run tests with additional options, you can use the following flags:

- `-s`: Show print statements in the test functions.
- `-v`: Show verbose output.
- `-k <expression>`: Run only tests that match the provided expression.
- `--collect-only`: List all test cases.
- `--image <image>`: Specify the Docker image to use for testing. If not provided, the `ubuntu2204` image is used. To use all available images, specify the `all` option.
- `--remove_images`: Remove the Docker images before running the tests to ensure a clean environment.

Example:

```bash
pytest -s -v -k "test_function_name"
```

To run all tests with all possible options, use the following command (running it may take hours):

```bash
pytest -s -v --image all --remove_images
```

To list all test cases, use the following command:

```bash
pytest --collect-only
```

## Packager tests

### 1. Test suite
