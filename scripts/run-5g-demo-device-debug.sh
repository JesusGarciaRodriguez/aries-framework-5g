#!/bin/bash
#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
set -e

BIN_PATH="${BIN_PATH}"
TLS_CERTS="$(cd ../test/bdd/fixtures/keys/tls && pwd)"
echo "$(cd ../test/bdd/fixtures/keys/tls && pwd)"

set -o allexport
set TLS_CERT_FILE=${TLS_CERTS}/ec-pubCert.pem
set TLS_KEY_FILE=${TLS_CERTS}/ec-key.pem
[[ -f 5g-device-poc/device_config.env ]] && source 5g-device-poc/device_config.env
# .env is missing the following TLS_CERT_FILE?=../test/bdd/fixtures/keys/tls/ec-pubCert.pem
# TLS_KEY_FILE?=../test/bdd/fixtures/keys/tls/ec-key.pem

cd $BIN_PATH 
echo ${ARIESD_API_HOST}
dlv --listen=:3000 --headless=true --api-version=2 --log=true --log-output=debugger,debuglineerr,gdbwire,lldbout,rpc --accept-multiclient exec ./aries-agent-rest --continue -- start
set +o allexport