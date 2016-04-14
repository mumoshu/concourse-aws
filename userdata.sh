#!/bin/bash -v
apt-get update -y
apt-get install -y nginx > /tmp/nginx.log
curl -v -L https://github.com/concourse/concourse/releases/download/v1.0.0/concourse_linux_amd64 -o concourse
chmod +x concourse
mv concourse /usr/local/bin/concourse
