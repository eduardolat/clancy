#!/bin/sh
set -e

# Usage:
# curl -sfL https://raw.githubusercontent.com/eduardolat/clancy/main/install.sh | sh

# Dependency Check
# Verify availability of essential tools before proceeding
DEPS="curl uname mktemp mv chmod"
for cmd in $DEPS; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Error: Required command '$cmd' not found."
    exit 1
  fi
done

# Configuration
OWNER="eduardolat"
REPO="clancy"
BIN_NAME="clancy"
INSTALL_DIR="/usr/local/bin"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize OS
case "$OS" in
  linux)  ;;
  darwin) ;;
  mingw*|msys*) OS="windows" ;;
  *) echo "Error: OS '$OS' is not supported."; exit 1 ;;
esac

# Normalize Architecture
case "$ARCH" in
  x86_64)        ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Error: Architecture '$ARCH' is not supported."; exit 1 ;;
esac

# Construct Asset Name
EXT=""
[ "$OS" = "windows" ] && EXT=".exe"
ASSET="${BIN_NAME}-${OS}-${ARCH}${EXT}"

URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${ASSET}"

# Execution
echo "Detected platform: $OS/$ARCH"
echo "Downloading latest release..."

TMP_FILE=$(mktemp)

# Download
if ! curl -sSfL "$URL" -o "$TMP_FILE"; then
  echo "Error: Failed to download asset. Please verify the release exists."
  rm -f "$TMP_FILE"
  exit 1
fi

chmod +x "$TMP_FILE"

# Installation
echo "Installing to $INSTALL_DIR..."
TARGET="$INSTALL_DIR/$BIN_NAME$EXT"

# Create directory if missing
if [ ! -d "$INSTALL_DIR" ]; then
  if [ -w "$(dirname "$INSTALL_DIR")" ]; then
    mkdir -p "$INSTALL_DIR"
  else
    sudo mkdir -p "$INSTALL_DIR"
  fi
fi

# Move binary
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_FILE" "$TARGET"
else
  sudo mv "$TMP_FILE" "$TARGET"
fi

echo "Successfully installed ${BIN_NAME} to ${TARGET}"
echo "Please ensure ${INSTALL_DIR} is in your PATH"
