#!/bin/bash

exec > /var/log/01_start_concourse_web.log 2>&1
set -x

CONCOURSE_PATH=/var/lib/concourse

mkdir -p $CONCOURSE_PATH

echo "${session_signing_key}" > $CONCOURSE_PATH/session_signing_key
echo "${tsa_host_key}" > $CONCOURSE_PATH/tsa_host_key
echo "${tsa_authorized_keys}" > $CONCOURSE_PATH/tsa_authorized_keys
echo "${postgres_data_source}" > $CONCOURSE_PATH/postgres_data_source
echo "${external_url}" > $CONCOURSE_PATH/external_url
curl http://169.254.169.254/latest/meta-data/local-ipv4 > $CONCOURSE_PATH/peer_ip

cd $CONCOURSE_PATH

concourse web --session-signing-key session_signing_key --tsa-host-key tsa_host_key --tsa-authorized-keys tsa_authorized_keys --external-url $(cat external_url) --postgres-data-source $(cat postgres_data_source) --basic-auth-username foo --basic-auth-password bar 2>&1 > $CONCOURSE_PATH/concourse_web.log &

echo $! > $CONCOURSE_PATH/pid
