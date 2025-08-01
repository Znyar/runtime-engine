#!/bin/bash

BASE_DIR="/data"

mkdir -p "$BASE_DIR"
cd "$BASE_DIR" || exit

curl -L -o "go-1.24.5.tar.gz" "https://go.dev/dl/go1.24.5.linux-amd64.tar.gz"