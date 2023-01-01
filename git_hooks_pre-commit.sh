#!/bin/bash

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# https://www.shellcheck.net/wiki/SC2155
repo_root="$(git rev-parse --show-toplevel)"
readonly repo_root

cd "$repo_root"

golangci-lint run

go test ./...
