resource "google_redis_instance" "aggregator" {
  name           = "aggregator"
  memory_size_gb = 1
  tier           = "BASIC"

  location_id = "europe-west1-b"

  authorized_network = "demo-net"

  redis_version = "REDIS_3_2"
  display_name  = "aggregator-test"

  labels {
    my_key    = "my_val"
    other_key = "other_val"
  }
}

data "google_dns_managed_zone" "private" {
  name     = "private-dns-zone"
}

resource "google_dns_record_set" "redis" {
  name = "redis.${data.google_dns_managed_zone.private.dns_name}"
  type = "CNAME"
  ttl  = 300

  managed_zone = "${data.google_dns_managed_zone.private.name}"

  rrdatas = ["${google_redis_instance.aggregator.host}"]
}


