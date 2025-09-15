# The script must be started in the root directory of the project
set -e

# Parse script arguments
SYSTEM=$1
ARCHITECTURE=$2
MODULE_DIR=$3

if [[ -z $SYSTEM || -z $ARCHITECTURE || -z $MODULE_DIR ]]; then
  echo "Usage: $0 <system> <architecture> <workdir>" >&2
  exit 1
fi

cd "$MODULE_DIR"

VERSION=$(sed -E -n 's/version=([^=]+)/\1/p' < version.txt)
MODULE_NAME=$(go list -m | awk -F/ '{print $NF}')
PACKAGE_NAME=${MODULE_NAME}_v${VERSION}_${SYSTEM}-${ARCHITECTURE}
INSTALL_DIR="./${PACKAGE_NAME}/"

if [ -d "$INSTALL_DIR" ]; then
  echo "Directory $INSTALL_DIR already exists. Please remove it before running the script." >&2
  exit 1
fi
if [ -f "${PACKAGE_NAME}.zip" ]; then
  echo "File ${PACKAGE_NAME}.zip already exists. Please remove it before running the script." >&2
  exit 1
fi

case $ARCHITECTURE in
  x86-64|x86)
    GO_ARCH="amd64"
    ;;
  aarch64)
    GO_ARCH="arm64"
    ;;
  aarch32)
    GO_ARCH="arm"
    ;;
  *)
    GO_ARCH=${ARCHITECTURE}
    ;;
esac

CGO_ENABLED=0 GOOS=${SYSTEM} GOARCH=${GO_ARCH} go build -a -o "${INSTALL_DIR}/bin/" ./...

if [ -d "resources" ]; then
  mkdir -p "$INSTALL_DIR/etc/"
	cp -r resources/* "$INSTALL_DIR/etc/"
fi

if [ -d "doc" ]; then
	cp -r doc "$INSTALL_DIR"
fi

cp README.md "$INSTALL_DIR"
cp version.txt "$INSTALL_DIR"
cp LICENSE "$INSTALL_DIR"

zip -r "${PACKAGE_NAME}.zip" "$INSTALL_DIR"
rm -rf "$INSTALL_DIR"