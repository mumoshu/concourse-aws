output "concourse_web_endpoint" {
  value = "${template_file.start_concourse_web.vars.external_url}"
}

output "concourse_web_elb_dns_name" {
  value = "${aws_elb.web-elb.dns_name}"
}
