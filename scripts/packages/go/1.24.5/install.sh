#!/bin/bash

INSTALL_DIR=/data/go/1.24.5
GO_ARCHIVE="/data/go-1.24.5.tar.gz"

mkdir -p "$INSTALL_DIR"
tar -xzf "$GO_ARCHIVE" -C "$INSTALL_DIR" --strip-components=1