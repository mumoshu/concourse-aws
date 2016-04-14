#!/bin/bash

packer build -var source_ami=$(./latest-ami-ubuntu.sh) docker-baked.json && packer build -var source_ami=$(./latest-ami-docker.sh) concourse-baked.json && ./my-latest-ami.sh && ./recreate.sh web
