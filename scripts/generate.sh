#!/bin/bash
OS_NAME=$(uname)
case $OS_NAME in
  Linux)
    BINARY_NAME=csvloader
    ;;
  Darwin)
    BINARY_NAME=csvloader
    ;;
  MINGW*|MSYS*|CYGWIN*)
    BINARY_NAME=csvloader.exe
    ;;
  *)
esac

BUILD_DIR=bin
SRC_DIR=cmd/csvloader
TEST_DIR=test
PKG_NAME=$SRC_DIR/csvloader.go
GO_BUILD_FLAGS=$BUILD_DIR/$BINARY_NAME
GO_TEST_FLAGS=-cover

build() {
  echo "Building $BINARY_NAME..."
  mkdir -p $BUILD_DIR
  go build -o $GO_BUILD_FLAGS $PKG_NAME
}

test() {
  echo "Running tests..."
  go test $GO_TEST_FLAGS $PKG_NAME
}

clean() {
  echo "Cleaning up..."
  rm -rf $BUILD_DIR
}

clean_all() {
  clean
  echo "Cleaning up all..."
  rm -rf $BUILD_DIR $TEST_DIR
}

help() {
  echo "Available commands:"
  echo "  build    - Build the application"
  echo "  test     - Run tests"
  echo "  clean    - Clean up build artifacts"
  echo "  clean-all- Clean up all artifacts including tests"
}

# Main function to parse arguments and call functions
main() {
  case "$1" in
    build)
      build
      ;;
    test)
      test
      ;;
    clean)
      clean
      ;;
    clean-all)
      clean_all
      ;;
    help)
      help
      ;;
    *)
      echo "Unknown command: $1"
      help
      ;;
  esac
}

main "$1"
