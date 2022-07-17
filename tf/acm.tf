#resource "aws_acm_certificate" "lb_listener" {
#  domain_name = aws_lb.main.dns_name
#  validation_method = "DNS"
#}