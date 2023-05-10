#!/usr/bin/env bash
set -xeou pipefail


DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)"

kubectl -n stackrox port-forward deploy/sensor 8443:8443 > /dev/null &
PID=$!
function ctrl_c() {
    echo "exit trap"
    rm -f "${DIR}/../certs/ca.pem" "${DIR}/../certs/collector-cert.pem" "${DIR}/../certs/collector-key.pem"
    kill $PID
}
trap ctrl_c INT

# -o=json query or
kubectl -n stackrox get secrets collector-tls -ojsonpath='{.data.ca\.pem}' | base64 -D -o "${DIR}/../certs/ca.pem"
kubectl -n stackrox get secrets collector-tls -ojsonpath='{.data.collector-cert\.pem}' | base64 -D -o "${DIR}/../certs/collector-cert.pem"
kubectl -n stackrox get secrets collector-tls -ojsonpath='{.data.collector-key\.pem}' | base64 -D -o "${DIR}/../certs/collector-key.pem"

ROX_MTLS_CA_FILE="${DIR}/../certs/ca.pem" \
ROX_MTLS_CERT_FILE="${DIR}/../certs/collector-cert.pem" \
ROX_MTLS_KEY_FILE="${DIR}/../certs/collector-key.pem" \
ROX_ADVERTISED_ENDPOINT="localhost:8443" \
go run tools/local-compliance/*.go
