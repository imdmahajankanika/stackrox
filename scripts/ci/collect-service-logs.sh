#!/bin/sh
set -eu

# Collect Service Logs script
#
# Extracts service logs from the given Kubernetes cluster and saves them off for
# future examination.
#
# Usage:
#   collect-service-logs.sh NAMESPACE [DIR]
#
# Example:
# $ ./scripts/ci/collect-service-logs.sh stackrox
#
# Assumptions:
# - Must be called from the root of the Apollo git repository.
# - Logs are saved under /tmp/k8s-service-logs/ or DIR if passed

usage() {
    echo "./scripts/ci/collect-service-logs.sh <namespace>"
    echo "e.g. ./scripts/ci/collect-service-logs.sh stackrox"
}

main() {
    namespace="$1"
    if [ -z "${namespace}" ]; then
        usage
        exit 1
    fi

    if [ $# -gt 1 ]; then
        log_dir="$2"
    else
        log_dir="/tmp/k8s-service-logs"
    fi
    log_dir="${log_dir}/${namespace}"
    mkdir -p "${log_dir}"

	set +e

    for object in deployments services pods secrets serviceaccounts; do
        # A feel good command before pulling logs
        echo ">>> ${object} <<<"
        kubectl -n "${namespace}" get "${object}" -o wide

        mkdir -p "${log_dir}/${object}"

        for item in $(kubectl -n "${namespace}" get "${object}" | tail -n +2 | awk '{print $1}'); do
            kubectl describe "${object}" "${item}" -n "${namespace}" > "${log_dir}/${object}/${item}_describe.log"
            for ctr in $(kubectl -n "${namespace}" get "${object}" "${item}" -o jsonpath='{.status.containerStatuses[*].name}'); do
                kubectl -n "${namespace}" logs "${object}/${item}" -c "${ctr}" > "${log_dir}/${object}/${item}-${ctr}.log"
                kubectl -n "${namespace}" logs "${object}/${item}" -p -c "${ctr}" > "${log_dir}/${object}/${item}-${ctr}-previous.log"
            done
        done
    done

    kubectl -n "${namespace}" get events -o wide >"${log_dir}/events.txt"
    find "${log_dir}" -type f -size 0 -delete
}

main "$@"
