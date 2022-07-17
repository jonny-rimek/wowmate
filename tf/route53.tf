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

resource "aws_route53_record" "validate_acm" {
  name    = "_49d84a4d2d415e0e9832b77ba5f71b91.test.dev.wowmate.io."
  type    = "CNAME"
  zone_id = "Z09026202SZR8MRVSF1BQ"
  records = ["_b9b15ea3e821425558ca92cee6446cb5.bwfqbhlrkg.acm-validations.aws."]
  ttl     = 300
}