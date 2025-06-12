output "nitric" {
  value = {
    id          = aws_lambda_function.function.arn
    domain_name = data.aws_lb.alb.dns_name
    exports = {
      resources = {
        "aws_lb" = var.alb_arn
        # The target port that this service has attached a listener for
        "aws_lb:target_port" = var.container_port
      }
    }
  }
}

