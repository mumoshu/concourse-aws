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

variable "web_instance_type" {
  default = "t2.micro"
  description = "AWS instance type for web"
}

variable "worker_instance_type" {
  default = "t2.micro"
  description = "AWS instance type for worker"
}

variable "asg_min" {
  description = "Min numbers of servers in ASG"
  default = "0"
}

variable "asg_max" {
  description = "Max numbers of servers in ASG"
  default = "2"
}

variable "web_asg_desired" {
  description = "Desired numbers of web servers in ASG"
  # Setting this gte 2 result in `fly execute --input foo=bar` to fail with errors like: "bad response uploading bits (404 Not Found)" or "gunzip: invalid magic"
  default = "1"
}

variable "worker_asg_desired" {
  description = "Desired numbers of servers in ASG"
  default = "2"
}

variable "elb_listener_lb_port" {
  description = ""
  default = "80"
}

variable "use_custom_elb_port" {
  default = 0
}

variable "elb_listener_lb_protocol" {
  default = "http"
}

variable "elb_listener_instance_port" {
  description = ""
  default = "8080"
}

variable "in_access_allowed_cidrs" {
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

variable "worker_instance_profile" {
  description = "IAM instance profile name to be used by Concourse workers. Can be an empty string to not specify it (no instance profile is used then)"
}

variable "basic_auth_username" {
  default = ""
}

variable "basic_auth_password" {
  default = ""
}

variable "github_auth_client_id" {
  default = ""
}

variable "github_auth_client_secret" {
  default = ""
}

variable "github_auth_organizations" {
  default = ""
}

variable "github_auth_teams" {
  default = ""
}

variable "github_auth_users" {
  default = ""
}

variable "custom_external_domain_name" {
  default = ""
  description ="don't include http[s]://"
}

variable "use_custom_external_domain_name" {
  default = 0
}

variable "ssl_certificate_arn" {
  default = ""
}
