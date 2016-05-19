# This article may help you understand what we do here
# https://dzone.com/articles/graceful-shutdown-using-aws-autoscaling-groups-and

resource "aws_sqs_queue" "graceful_termination_queue" {
  name = "graceful_termination_queue"
}

resource "aws_iam_role" "autoscaling_role" {
  name = "autoscaling_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "autoscaling.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "lifecycle_hook_autoscaling_policy" {
  name = "lifecycle_hook_autoscaling_policy"
  role = "${aws_iam_role.autoscaling_role.id}"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1436380187000",
            "Effect": "Allow",
            "Action": [
                "sqs:GetQueueUrl",
                "sqs:SendMessage"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
EOF
}

resource "aws_autoscaling_lifecycle_hook" "graceful_shutdown_asg_hook" {
  name = "graceful_shutdown_asg"
  autoscaling_group_name = "${var.target_asg_name}"
  # When a hook is timed out or failed unexpectedly,
  # we want not to ABANDON but to CONTINUE the remaining auto scaling process.
  default_result = "CONTINUE"
  heartbeat_timeout = 60
  lifecycle_transition = "autoscaling:EC2_INSTANCE_TERMINATING"
  notification_target_arn = "${aws_sqs_queue.graceful_termination_queue.arn}"
  role_arn = "${aws_iam_role.autoscaling_role.arn}"
}

output "sqs_queue_arn" {
  value = "${aws_sqs_queue.graceful_termination_queue.arn}"
}

variable "target_asg_name" {
}
