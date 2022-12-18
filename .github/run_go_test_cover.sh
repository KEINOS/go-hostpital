#!/bin/sh
# =============================================================================
#  This script runs the tests with coverage. It will fail if the coverage is
#  below 100.0%.
# =============================================================================

set -eu

go test -cover ./... | grep ok | grep 100 || {
    echo 'ERROR: coverage is below 100.0%'
    exit 1
}