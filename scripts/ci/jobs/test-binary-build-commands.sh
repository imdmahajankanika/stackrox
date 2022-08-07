#!/usr/bin/env bash

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")"/../../.. && pwd)"
source "$ROOT/scripts/ci/lib.sh"

set -euo pipefail

make_test_bin() {
    info "Making test-bin"

    export OPENSHIFT_BUILD_NAME="test-bin"

    make cli-build upgrader
    install_built_roxctl_in_gopath
}

make_test_bin "$*"
