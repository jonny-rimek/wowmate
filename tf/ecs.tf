resource "aws_ecs_cluster" "main" {
  name = "test-cluster"
}

resource "aws_ecs_task_definition" "main" {
  network_mode             = "awsvpc"
  family                   = "service"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  container_definitions    = jsonencode([
    {
      name         = "test-container"
      image        = "nginx:latest"
      essential    = true
      portMappings = [
        {
          protocol      = "tcp"
          containerPort = var.container_port
          hostPort      = var.container_port
        }
      ]
    }
  ])
}

resource "aws_ecs_service" "main" {
  name                               = "test-service"
  cluster                            = aws_ecs_cluster.main.id
  task_definition                    = aws_ecs_task_definition.main.arn
  desired_count                      = var.number_of_containers
  deployment_minimum_healthy_percent = 50
  deployment_maximum_percent         = 200
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"

  network_configuration {
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets          = aws_subnet.private.*.id
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.main.arn
    container_name   = "test-container"
    container_port   = var.container_port
  }

#  lifecycle {
#    ignore_changes = [task_definition, desired_count]
#  }
}