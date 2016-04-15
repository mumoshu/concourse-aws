#!/bin/bash

set -x

if [ ! -e host_key ]; then
  ssh-keygen -t rsa -f host_key -N ''
fi

if [ ! -e worker_key ]; then
  ssh-keygen -t rsa -f worker_key -N ''
fi

if [ ! -e session_signing_key ]; then
  ssh-keygen -t rsa -f session_signing_key -N ''
fi

cp worker_key.pub authorized_worker_keys

subnet_id=$CONCOURSE_SUBNET_ID

vpc_id=$(aws ec2 describe-subnets --subnet-id $subnet_id | jq -r .Subnets[].VpcId)

echo $vpc_id

terraform "$@" -var aws_region=ap-northeast-1 -var availability_zones=ap-northeast-1c -var key_name=cw_kuoka -var subnet_id=$subnet_id -var vpc_id=$vpc_id -var db_instance_class=db.t2.micro -var db_username=concourse -var db_password=concourse -var db_subnet_ids=$CONCOURSE_DB_SUBNET_IDS \
  -var tsa_host_key=host_key \
  -var session_signing_key=session_signing_key \
  -var tsa_authorized_keys=worker_key.pub \
  -var tsa_public_key=host_key.pub \
  -var tsa_worker_private_key=worker_key \
  -var ami=$(./my-latest-ami.sh) \
  -var in_access_allowed_cidr=$CONCOURSE_IN_ACCESS_ALLOWED_CIDR \
  -var worker_instance_profile=$CONCOURSE_WORKER_INSTANCE_PROFILE
