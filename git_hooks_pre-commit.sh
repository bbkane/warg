#!/bin/bash

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# https://www.shellcheck.net/wiki/SC2155
# https://stackoverflow.com/a/957978/2958070
repo_root="$(git rev-parse --show-toplevel)"
readonly repo_root

cd "$repo_root"

# an arg might not be passed to the script and if so, then "$1" will be unset.
# Temporarily disable unset error checking for to account for this
set +u
if [ "$1" == "link" ]; then
    set -x
    ln -s "${repo_root}/git_hooks_pre-commit.sh" "${repo_root}/.git/hooks/pre-commit"
    { set +x; } 2>/dev/null
    exit 0
elif [ "$1" == "unlink" ]; then
    set -x
    unlink "${repo_root}/.git/hooks/pre-commit"
    { set +x; } 2>/dev/null
    exit 0
fi
set -u

golangci-lint run

go test ./...
