resource "aws_ecs_task_definition" "api" {
  family                   = "${local.prefix}-api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.task_execution_role.arn
  task_role_arn            = aws_iam_role.app_task.arn

  container_definitions = jsonencode(
    [
      {
        name              = "api"
        image             = var.ecr_app_image
        essential         = true
        memoryReservation = 256
        portMappings = [
          {
            containerPort = 8080
            hostPort      = 8080
          }
        ]
        environment = [
          {
            name  = "AWASM_DB_HOST"
            value = aws_db_instance.main.address
          },
          {
            name  = "AWASM_DB_NAME"
            value = aws_db_instance.main.db_name
          },
          {
            name  = "AWASM_DB_USER"
            value = aws_db_instance.main.username
          },
          {
            name  = "AWASM_DB_PASSWORD"
            value = aws_db_instance.main.password
          },
          {
            name  = "AWASM_DB_SCHEMA"
            value = aws_db_instance.main.db_name
          },
          {
            name  = "AWASM_REDIS_HOST"
            value = aws_elasticache_cluster.redis.cache_nodes[0].address
          },
          {
            name  = "AWASM_MAILER_PROVIDER_TYPE"
            value = "NOOP"
          },
        ]
        logConfiguration = {
          logDriver = "awslogs"
          options = {
            awslogs-group         = aws_cloudwatch_log_group.ecs_task_logs.name
            awslogs-region        = data.aws_region.current.name
            awslogs-stream-prefix = "api"
          }
        }
      }
    ]
  )

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }
}

resource "aws_ecs_service" "api" {
  name                   = "${local.prefix}-api"
  cluster                = aws_ecs_cluster.main.name
  task_definition        = aws_ecs_task_definition.api.family
  desired_count          = 1
  launch_type            = "FARGATE"
  platform_version       = "1.4.0"
  enable_execute_command = true

  network_configuration {
    # TOIMPROVE: in production, we use vpc endpoint to improve security
    assign_public_ip = true

    subnets = [
      aws_subnet.public_a.id,
      aws_subnet.public_b.id,
      aws_subnet.private_a.id,
      aws_subnet.private_b.id,
    ]

    security_groups = [aws_security_group.ecs_service.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api"
    container_port   = 8080
  }
}
