#!/bin/bash

aws ec2 describe-images --owners self --filters Name=virtualization-type,Values=hvm Name=root-device-type,Values=ebs Name=architecture,Values=x86_64 Name=name,Values="packer-ubuntu-xenial-docker-*" | jq -r ".Images | sort_by(.CreationDate) | .[].ImageId" | tail -n 1
