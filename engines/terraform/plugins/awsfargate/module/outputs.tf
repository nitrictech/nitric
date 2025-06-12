output "nitric" {
  value = {
    id          = aws_ecs_service.service.id
    domain_name = data.aws_lb.alb.dns_name
    exports = {
      resources = {
        "aws_lb" = var.alb_arn
        # The target port that this service has attached a listener for
        "aws_lb:http_port" = var.container_port
        # TODO: Fix the HTTPS port assignment
        "aws_lb:https_port" = var.container_port + 1
      }
    }
  }
}

