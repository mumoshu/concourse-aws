#!/bin/bash

function web() {
  terraform taint template_cloudinit_config.web && terraform taint aws_autoscaling_group.web-asg && ./apply.sh apply
}

function worker() {
  terraform taint template_cloudinit_config.worker && terraform taint aws_autoscaling_group.worker-asg && ./apply.sh apply
}

"$@"
