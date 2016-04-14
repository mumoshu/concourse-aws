# Specify the provider and access details
provider "aws" {
  region = "${var.aws_region}"
}

#module "postgres" {
#    source = "./postgres"
#    access_allowed_security_groups = "${aws_security_group.atc.id}"
#}

resource "aws_elb" "web-elb" {
  name = "terraform-example-elb"

  # The same availability zone as our instances
  # Only one of SubnetIds or AvailabilityZones may be specified
  #availability_zones = ["${split(",", var.availability_zones)}"]
  security_groups = ["${aws_security_group.external_lb.id}"]
  subnets = ["${var.subnet_id}"]

  listener {
    instance_port = "${var.elb_listener_instance_port}"
    instance_protocol = "http"
    lb_port = "${var.elb_listener_lb_port}"
    lb_protocol = "http"
  }

  listener {
    instance_port = "${var.tsa_port}"
    instance_protocol = "tcp"
    lb_port = "${var.tsa_port}"
    lb_protocol = "tcp"
  }

  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "TCP:${var.elb_listener_instance_port}"
    interval = 30
  }

}

resource "aws_autoscaling_group" "web-asg" {
  availability_zones = ["${split(",", var.availability_zones)}"]
  name = "terraform-example-asg"
  max_size = "${var.asg_max}"
  min_size = "${var.asg_min}"
  desired_capacity = "${var.asg_desired}"
  force_delete = true
  launch_configuration = "${aws_launch_configuration.web-lc.name}"
  load_balancers = ["${aws_elb.web-elb.name}"]
  vpc_zone_identifier = ["${split(",", var.subnet_id)}"]
  tag {
    key = "Name"
    value = "${var.prefix}web"
    propagate_at_launch = "true"
  }
}

resource "aws_autoscaling_group" "worker-asg" {
  availability_zones = ["${split(",", var.availability_zones)}"]
  max_size = "${var.asg_max}"
  min_size = "${var.asg_min}"
  desired_capacity = "${var.asg_desired}"
  force_delete = true
  launch_configuration = "${aws_launch_configuration.worker-lc.name}"
  vpc_zone_identifier = ["${split(",", var.subnet_id)}"]
  tag {
    key = "Name"
    value = "${var.prefix}worker"
    propagate_at_launch = "true"
  }
}

resource "aws_launch_configuration" "web-lc" {
  # Omit launch configuration name to avoid collisions on create_before_destroy
  # ref. https://github.com/hashicorp/terraform/issues/1109#issuecomment-97970885
  #image_id = "${lookup(var.aws_amis, var.aws_region)}"
  image_id = "${var.ami}"
  instance_type = "${var.instance_type}"
  security_groups = ["${aws_security_group.default.id}","${aws_security_group.atc.id}","${aws_security_group.tsa.id}"]
  user_data = "${template_cloudinit_config.web.rendered}"
  key_name = "${var.key_name}"
  associate_public_ip_address = true
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_launch_configuration" "worker-lc" {
  #image_id = "${lookup(var.aws_amis, var.aws_region)}"
  image_id = "${var.ami}"
  instance_type = "${var.instance_type}"
  security_groups = ["${aws_security_group.default.id}", "${aws_security_group.worker.id}"]
  user_data = "${template_cloudinit_config.worker.rendered}"
  key_name = "${var.key_name}"
  associate_public_ip_address = true
  lifecycle {
    create_before_destroy = true
  }
}
  
resource "template_file" "install_concourse" {
  template = "${file("${path.module}/00_install_concourse.sh.tpl")}"

#  vars {
#    consul_address = "${aws_instance.consul.private_ip}"
#  }
}

resource "template_file" "start_concourse_web" {
  template = "${file("${path.module}/01_start_concourse_web.sh.tpl")}"

  vars {
    session_signing_key = "${file("${path.module}/${var.session_signing_key}")}"
    tsa_host_key = "${file("${path.module}/${var.tsa_host_key}")}"
    tsa_authorized_keys = "${file("${path.module}/${var.tsa_authorized_keys}")}"
    postgres_data_source = "postgres://${var.db_username}:${var.db_password}@${aws_db_instance.default.endpoint}/concourse"
    external_url = "http://${aws_elb.web-elb.dns_name}"
    # peer_url
  }
}

resource "template_file" "start_concourse_worker" {
  template = "${file("${path.module}/02_start_concourse_worker.sh.tpl")}"

  vars {
    tsa_host = "${aws_elb.web-elb.dns_name}"
    tsa_public_key = "${file("${path.module}/${var.tsa_public_key}")}"
    tsa_worker_private_key = "${file("${path.module}/${var.tsa_worker_private_key}")}"
  }
}

resource "template_cloudinit_config" "web" {
  gzip          = false
  base64_encode = false

  part {
    content_type = "text/x-shellscript"
    content      = "${template_file.install_concourse.rendered}"
  }

  part {
    content_type = "text/x-shellscript"
    content      = "${template_file.start_concourse_web.rendered}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "template_cloudinit_config" "worker" {
  gzip          = false
  base64_encode = false

  part {
    content_type = "text/x-shellscript"
    content      = "${template_file.install_concourse.rendered}"
  }

  part {
    content_type = "text/x-shellscript"
    content      = "${template_file.start_concourse_worker.rendered}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "default" {
  name_prefix = "${var.prefix}default"
  description = "Used in the terraform"
  vpc_id = "${var.vpc_id}"

  # SSH access from anywhere
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["${var.in_access_allowed_cidr}"]
  }

  # outbound internet access
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "atc" {
  name_prefix = "${var.prefix}atc"
  description = "Used in the terraform"
  vpc_id = "${var.vpc_id}"

  # HTTP access from anywhere
  ingress {
    from_port = "${var.elb_listener_instance_port}"
    to_port = "${var.elb_listener_instance_port}"
    protocol = "tcp"
    cidr_blocks = ["${var.in_access_allowed_cidr}"]
  }
}

resource "aws_security_group_rule" "allow_external_lb_to_atc_access" {
    type = "ingress"
    from_port = "${var.elb_listener_instance_port}"
    to_port = "${var.elb_listener_instance_port}"
    protocol = "tcp"

    security_group_id = "${aws_security_group.tsa.id}"
    source_security_group_id = "${aws_security_group.external_lb.id}"
}

resource "aws_security_group_rule" "allow_atc_to_worker_access" {
    type = "ingress"
    from_port = "0"
    to_port = "65535"
    protocol = "tcp"

    security_group_id = "${aws_security_group.worker.id}"
    source_security_group_id = "${aws_security_group.atc.id}"
}

resource "aws_security_group" "tsa" {
  name_prefix = "${var.prefix}tsa"
  description = "Used for concourse ci tsa"
  vpc_id = "${var.vpc_id}"

  # outbound internet access
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group_rule" "allow_worker_to_tsa_access" {
    type = "ingress"
    from_port = 2222
    to_port = 2222
    protocol = "tcp"

    security_group_id = "${aws_security_group.tsa.id}"
    source_security_group_id = "${aws_security_group.worker.id}"
}

resource "aws_security_group_rule" "allow_external_lb_to_tsa_access" {
    type = "ingress"
    from_port = 2222
    to_port = 2222
    protocol = "tcp"

    security_group_id = "${aws_security_group.tsa.id}"
    source_security_group_id = "${aws_security_group.external_lb.id}"
}

resource "aws_security_group" "worker" {
  name_prefix = "${var.prefix}worker"
  description = "Used for concourse ci worker"
  vpc_id = "${var.vpc_id}"

  # outbound internet access
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "external_lb" {
  #name = "external_lb"
  description = "Used in the terraform"

  vpc_id = "${var.vpc_id}"

  ingress {
    from_port = "${var.elb_listener_lb_port}"
    to_port = "${var.elb_listener_lb_port}"
    protocol = "tcp"
    cidr_blocks = ["${var.in_access_allowed_cidr}"]
  }

  ingress {
    from_port = "${var.tsa_port}"
    to_port = "${var.tsa_port}"
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "db" {
  name_prefix = "${var.prefix}db"
  description = "Used for concourse ci db"
  vpc_id = "${var.vpc_id}"

  # outbound internet access
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group_rule" "allow_db_access_from_atc" {
    type = "ingress"
    from_port = 5432
    to_port = 5432
    protocol = "tcp"

    security_group_id = "${aws_security_group.db.id}"
    source_security_group_id = "${aws_security_group.atc.id}"
}

resource "aws_db_instance" "default" {
  depends_on = ["aws_security_group.db"]
  identifier = "${var.prefix}db"
  allocated_storage = "10"
  engine = "postgres"
  engine_version = "9.4.1"
  instance_class = "${var.db_instance_class}"
  name = "concourse"
  username = "${var.db_username}"
  password = "${var.db_password}"
  vpc_security_group_ids = ["${aws_security_group.db.id}"]
  db_subnet_group_name = "${aws_db_subnet_group.db.id}"
}

resource "aws_db_subnet_group" "db" {
  name = "${var.prefix}db"
  description = "group of subnets for concourse db"
  subnet_ids = ["${split(",", var.db_subnet_ids)}"]
}
