#!/bin/bash

set -eu

packer build -var source_ami=ami-5d38d93c docker-baked.json
