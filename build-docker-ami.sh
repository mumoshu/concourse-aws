#!/bin/bash

set -eu

packer build -var source_ami=ami-0417e362 docker-baked.json
