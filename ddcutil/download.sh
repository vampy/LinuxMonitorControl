#!/usr/bin/env bash
set -euf -o pipefail
SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

DOWNLOAD_LINK="http://www.ddcutil.com/tarballs/ddcutil-1.0.1.tar.gz"
FILENAME="$SCRIPT_DIR/ddcutil.tar.gz"
DESTINATION_DIR="$SCRIPT_DIR/src"

echo -e "\n=== Downloading === \n"
wget --verbose -O "$FILENAME" "$DOWNLOAD_LINK"

echo -e "\n=== Unarchiving === \n"

echo "Removing $DESTINATION_DIR"
rm -rf "$DESTINATION_DIR"

echo "Creating $DESTINATION_DIR"
mkdir -p "$DESTINATION_DIR"

echo "Unarchiving $FILENAME into $DESTINATION_DIR"
tar xzf "$FILENAME" -C "$DESTINATION_DIR" --strip-components 1

echo -e "\n=== Cleaning up === \n"
rm --verbose "$FILENAME"
