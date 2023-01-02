#!/bin/bash

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

make_print_color() {
    color_name="$1"
    color_code="$2"
    color_reset="$(tput sgr0)"
    if [ -t 1 ] ; then
        eval "print_${color_name}() { printf \"${color_code}%s${color_reset}\\n\" \"\$1\"; }"
    else  # Don't print colors on pipes
        eval "print_${color_name}() { printf \"%s\\n\" \"\$1\"; }"
    fi
}

# https://unix.stackexchange.com/a/269085/185953
make_print_color "red" "$(tput setaf 1)"
make_print_color "green" "$(tput setaf 2)"
make_print_color "yellow" "$(tput setaf 3)"
make_print_color "blue" "$(tput setaf 4)"

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

print_blue "# Running: pre-commit"

print_blue "## Running: golangci-lint"
golangci-lint run || { print_red "Failed!"; exit 1; }
print_green "## Succeeded: golangci-lint"

print_blue "## Running: go test"
go test ./... > /dev/null || { print_red "Failed!"; exit 1; }
print_green "## Succeeded: go test"

print_green "# Succeeded: pre-commit"
