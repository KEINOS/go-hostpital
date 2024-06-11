#!/bin/sh
# =============================================================================
#  This script updates Go modules to the latest version.
# =============================================================================
#  It will backup the go.mod and go.sum files, then run `go get -u ./...` to
#  update the modules.
#  Then it will run the tests to make sure the code is still working, and fails
#  if any errors are found during the process. Fallback to the original in case
#  of errors.
#
#  NOTE: This script is aimed to run in the container via docker-compose.
#    See "tidy" service: ./docker-compose.yml
# =============================================================================

min_go_version='1.22'

set -eu

print_ok() {
    printf 'ok\t%s\n' "$1"
}

# -----------------------------------------------------------------------------
echo '* Backup module files ...'

cp go.mod go.mod.bak || {
    echo 'ERROR: failed to backup go.mod file'
    exit 1
}
cp go.sum go.sum.bak || {
    echo 'ERROR: failed to backup go.sum file'
    exit 1
}

print_ok 'go.mod.bak and go.sum.bak created'

# -----------------------------------------------------------------------------
echo '* Run go tidy ...'
go mod tidy -go=${min_go_version}|| {
    echo 'error: failed to run go mod tidy'
    echo '!!: Plese fallback to the original files'
    exit 1
}

print_ok 'go mod tidy done'

# -----------------------------------------------------------------------------
echo '* Run tests ...'
go test ./... || {
    echo 'ERROR: failed to run tests'
    exit 1
}

# -----------------------------------------------------------------------------
echo '* Removing old module files ...'
rm -f go.mod.bak || {
    echo 'ERROR: failed to remove go.mod.bak file'
    exit 1
}

rm -f go.sum.bak || {
    echo 'ERROR: failed to remove go.sum.bak file'
    exit 1
}

print_ok 'go.mod.bak and go.sum.bak removed'

# -----------------------------------------------------------------------------
echo
echo 'Modules successfully updated!'
