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
echo "${github_auth_organizations}" > $CONCOURSE_PATH/github_auth_arganizations
echo "${github_auth_teams}" > $CONCOURSE_PATH/github_auth_teams
echo "${github_auth_users}" > $CONCOURSE_PATH/github_auth_users
curl http://169.254.169.254/latest/meta-data/local-ipv4 > $CONCOURSE_PATH/peer_ip

if [ "z${basic_auth_username}" != "z" ]; then
  BASIC_AUTH_OPTS="--basic-auth-username ${basic_auth_username} --basic_auth_password ${basic_auth_password}"
fi

GITHUB_AUTH_OPTS=()
if [ "z${github_auth_client_id}" != "z" ]; then
  GITHUB_AUTH_OPTS+=("--github-auth-client-id")
  GITHUB_AUTH_OPTS+=("${github_auth_client_id}")
  GITHUB_AUTH_OPTS+=("--github-auth-client-secret")
  GITHUB_AUTH_OPTS+=("${github_auth_client_secret}")

  if [ "z${github_auth_organizations}" != "z" ]; then
    str="${github_auth_organizations}"
    IFS_ORIGINAL="$$IFS"
    IFS=,
    arr=($$str)
    IFS="$$IFS_ORIGINAL"
    for o in "$${arr[@]}"; do
      GITHUB_AUTH_OPTS+=("--github-auth-organization")
      GITHUB_AUTH_OPTS+=("$$o")
    done
  fi
  if [ "z${github_auth_teams}" != "z" ]; then
    str="${github_auth_teams}"
    IFS_ORIGINAL="$$IFS"
    IFS=,
    arr=($$str)
    IFS="$$IFS_ORIGINAL"
    for t in "$${arr[@]}"; do
      GITHUB_AUTH_OPTS+=("--github-auth-team")
      GITHUB_AUTH_OPTS+=("$$t")
    done
  fi
  if [ "z${github_auth_users}" != "z" ]; then
    str="${github_auth_users}"
    IFS_ORIGINAL="$$IFS"
    IFS=,
    arr=($$str)
    IFS="$$IFS_ORIGINAL"
    for u in "$${arr[@]}"; do
      GITHUB_AUTH_OPTS+=("--github-auth-user")
      GITHUB_AUTH_OPTS+=("$$u")
    done
  fi
fi

cd $CONCOURSE_PATH

concourse web --session-signing-key session_signing_key \
  --tsa-host-key tsa_host_key --tsa-authorized-keys tsa_authorized_keys \
  --external-url $(cat external_url) \
  --postgres-data-source $(cat postgres_data_source) \
  $BASIC_AUTH_OPTS \
  "$${GITHUB_AUTH_OPTS[@]}" \
  2>&1 > $CONCOURSE_PATH/concourse_web.log &

echo $! > $CONCOURSE_PATH/pid
