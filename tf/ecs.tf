resource "aws_ecs_cluster" "django_cluster" {
  name = "django"
}

resource "aws_ecs_task_definition" "django_task_definition" {
  container_definitions = jsonencode([
    {
      name      = "first"
      image     = "nginx"
      cpu       = 10
      memory    = 512
      essential = true
      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
      }]
  }])
  family = "service"
}

resource "aws_ecs_service" "django_service" {
  name = "django"
  cluster = aws_ecs_cluster.django_cluster.id
  task_definition = aws_ecs_task_definition.django_task_definition.arn
  desired_count = 1
}