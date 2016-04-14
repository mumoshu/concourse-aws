output "concourse_web_dns_name" {
  value = "${aws_elb.web-elb.dns_name}"
}
