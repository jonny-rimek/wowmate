resource "aws_route53_record" "test" {
  name    = "test.dev.wowmate.io"
  type    = "A"
  zone_id = "Z09026202SZR8MRVSF1BQ"
  
  alias {
    evaluate_target_health = false
    name                   = aws_lb.main.dns_name
    zone_id                = aws_lb.main.zone_id
  }
}