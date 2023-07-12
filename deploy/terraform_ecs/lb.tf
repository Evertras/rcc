resource "aws_security_group" "lb" {
  name   = "${local.prefix}-lb"
  vpc_id = aws_vpc.default.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "api" {
  name    = "${local.prefix}-api"
  subnets = aws_subnet.public.*.id

  security_groups = [aws_security_group.lb.id]
}

resource "aws_lb_target_group" "api" {
  name        = "${local.prefix}-api"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.default.id
  target_type = "ip"

  health_check {
    enabled = true
    // For now, until we have a proper health check endpoint
    matcher = "404"
  }
}

resource "aws_lb_listener" "api" {
  load_balancer_arn = aws_lb.api.id
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_lb_target_group.api.id
    type             = "forward"
  }
}
