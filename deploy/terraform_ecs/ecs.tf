resource "aws_ecs_task_definition" "api" {
  family                   = "${local.prefix}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]

  cpu    = 256
  memory = 512

  container_definitions = <<DEFINITION
[
  {
    "image": "evertras/rcc:v0.1.0",
    "cpu": 256,
    "memory": 512,
    "name": "${local.prefix}-api",
    "networkMode": "awsvpc",
    "portMappings": [
      {
        "containerPort": ${local.port},
        "hostPort": ${local.port}
      }
    ]
  }
]
DEFINITION
}

resource "aws_security_group" "api_task" {
  name   = "${local.prefix}-api"
  vpc_id = aws_vpc.default.id

  ingress {
    protocol        = "tcp"
    from_port       = local.port
    to_port         = local.port
    security_groups = [aws_security_group.lb.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_cluster" "main" {
  name = "${local.prefix}-main"
}

resource "aws_ecs_service" "api" {
  name            = "${local.prefix}-api"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api.arn
  desired_count   = var.app_count
  launch_type     = "FARGATE"

  network_configuration {
    security_groups = [aws_security_group.api_task.id]
    subnets         = aws_subnet.private.*.id
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.id
    container_name   = "${local.prefix}-api"
    container_port   = local.port
  }

  depends_on = [aws_lb_listener.api]
}
