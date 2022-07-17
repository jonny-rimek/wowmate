resource "aws_acm_certificate" "lb_listener" {
  domain_name = var.domain
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}