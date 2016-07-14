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
if [ "z${github_auth_client_id}" != "z" ]; then
  GITHUB_CLIENT="--github-auth-client-id ${github_auth_client_id} --github-auth-client-secret ${github_auth_client_secret}"

  if [ "z${github_auth_organizations}" != "z" ]; then
    GITHUB_AUTH_ORGS=$(echo "${github_auth_organizations}" | sed -e 's/^/--github-auth-organization /' -e 's/,/ --github-auth-organization /g')
  fi
  if [ "z${github_auth_teams}" != "z" ]; then
    GITHUB_AUTH_TEAMS=$(echo "${github_auth_teams}" | sed -e 's/^/--github-auth-team /' -e 's/,/ --github-auth-team /g')
  fi
  if [ "z${github_auth_users}" != "z" ]; then
    GITHUB_AUTH_USERS=$(echo "${github_auth_users}" | sed -e 's/^/--github-auth-user /' -e 's/,/ --github-auth-user /g')
  fi
  GITHUB_AUTH_OPTS="$GITHUB_CLIENT $GITHUB_AUTH_ORGS $GITHUB_AUTH_TEAMS $GITHUB_AUTH_USERS"
fi

cd $CONCOURSE_PATH

concourse web --session-signing-key session_signing_key \
  --tsa-host-key tsa_host_key --tsa-authorized-keys tsa_authorized_keys \
  --external-url $(cat external_url) \
  --postgres-data-source $(cat postgres_data_source) \
  $BASIC_AUTH_OPTS \
  $GITHUB_AUTH_OPTS \
  2>&1 > $CONCOURSE_PATH/concourse_web.log &

echo $! > $CONCOURSE_PATH/pid
