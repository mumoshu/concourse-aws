resource "aws_autoscaling_policy" "add-worker" {
    name = "${var.target_asg_name}-add-worker"
    scaling_adjustment = 1
    adjustment_type = "ChangeInCapacity"
    cooldown = 300
    autoscaling_group_name = "${var.target_asg_name}"
}

resource "aws_cloudwatch_metric_alarm" "worker-is-busy" {
    alarm_name = "${var.target_asg_name}-is-busy"
    comparison_operator = "GreaterThanOrEqualToThreshold"
    evaluation_periods = "2"
    metric_name = "CPUUtilization"
    namespace = "AWS/EC2"
    period = "120"
    statistic = "Average"
    threshold = "80"
    dimensions {
        AutoScalingGroupName = "${var.target_asg_name}"
    }
    alarm_description = "This metric monitor ec2 cpu utilization"
    alarm_actions = ["${aws_autoscaling_policy.add-worker.arn}"]
}

variable "target_asg_name" {}
