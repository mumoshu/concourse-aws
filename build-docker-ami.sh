#!/bin/bash

set -eu

packer build -var source_ami=ami-b7d829d6 docker-baked.json
