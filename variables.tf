variable "prefix" {
  description = "Prefix for every resource created by this template"
  default = "concourse-"  
}

variable "aws_region" {
  description = "The AWS region to create things in."
  default = "us-east-1"
}

# ubuntu-trusty-14.04 (x64)
variable "aws_amis" {
  default = {
    "us-east-1" = "ami-5f709f34"
    "us-west-2" = "ami-7f675e4f"
    "ap-northeast-1" = "ami-a21529cc"
  }
}

variable "ami" {
}

variable "availability_zones" {
  default = "us-east-1b,us-east-1c,us-east-1d,us-east-1e"
  description = "List of availability zones, use AWS CLI to find your "
}

variable "key_name" {
  description = "Name of AWS key pair"
}

variable "instance_type" {
  default = "t2.micro"
  description = "AWS instance type"
}

variable "asg_min" {
  description = "Min numbers of servers in ASG"
  default = "1"
}

variable "asg_max" {
  description = "Max numbers of servers in ASG"
  default = "2"
}

variable "asg_desired" {
  description = "Desired numbers of servers in ASG"
  default = "1"
}

variable "elb_listener_lb_port" {
  description = ""
  default = "80"
}

variable "elb_listener_instance_port" {
  description = ""
  default = "8080"
}

variable "in_access_allowed_cidr" {
  description = ""
}

variable "subnet_id" {
  description = ""
}

variable "db_subnet_ids" {
  description = ""
}

variable "vpc_id" {
  description = ""
}

variable "db_username" {
  description = ""
}

variable "db_password" {
  description = ""
}

variable "db_instance_class" {
  description = "t2.micro"
}

variable "tsa_host_key" {
  description = ""
}

variable "session_signing_key" {
  description = ""
}

variable "tsa_authorized_keys" {
  description = ""
}

variable "tsa_public_key" {
  description = ""
}

variable "tsa_worker_private_key" {
  description = ""
}

variable "tsa_port" {
  description = ""
  default = "2222"
}
