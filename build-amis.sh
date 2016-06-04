#!/bin/sh

set -eu

./build-docker-ami.sh
./build-concourse-ami.sh
