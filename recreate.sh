#!/bin/bash

function web() {
  terraform taint template_cloudinit_config.web && ./terraform.sh apply
}

function worker() {
  terraform taint template_cloudinit_config.worker && ./terraform.sh apply
}

"$@"
