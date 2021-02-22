#!/usr/bin/env bash
set -eu -o pipefail
SCRIPT_DIR="$(dirname "$(readlink -f "$0")")"

SRC_DIR="$SCRIPT_DIR/src"
DST_LIB_DIR="$SCRIPT_DIR/lib"
DST_BIN_DIR="$SCRIPT_DIR/bin"
DOWNLOAD_SCRIPT="$SCRIPT_DIR/download.sh"

if [[ ! -d "$SRC_DIR" ]]; then
    echo "Warning: $SRC_DIR does not exist"
    echo "Running Download script: $DOWNLOAD_SCRIPT"
    bash "$DOWNLOAD_SCRIPT"

    if [[ ! -d "$SRC_DIR" ]]; then
        exit "Error: $SRC_DIR  STILL does not exist"
    fi
fi


echo -e "\n=== Building === \n"

pushd "$SRC_DIR"

if [[ -d "src/.libs" ]]; then
    make clean
fi

rm -rf "$DST_BIN_DIR"
rm -rf "$DST_LIB_DIR"

# TODO https://github.com/rockowitz/ddcutil/issues/183
# --enable-static
./configure --disable-drm --prefix="$SCRIPT_DIR"
make -j 4
make install
popd


# echo -e "\n=== Copying === \n"
# mkdir -p "$DST_BIN_DIR"
# cp --verbose  "$SRC_DIR/src/ddcutil" "$DST_BIN_DIR/"

# Enable copying of static libs
# mkdir -p "$DST_LIB_DIR"
# cp --verbose --recursive --force "$SRC_DIR/src/.libs/." "$DST_LIB_DIR/"
