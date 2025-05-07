resource "aws_route53_record" "mst-domain" {
  zone_id = "Z04438833D2EQ4Q443N2K"
  name    = "metasoundtools.com"
  type    = "A"

  alias {
    name                   = "d2m8tnimif7e48.cloudfront.net"
    zone_id                = "Z2FDTNDATAQYW2"
    evaluate_target_health = false
  }
}
