output "load_balancer_dns" {
  value = aws_lb.main.dns_name
}