provider "google" {
  region  = "${var.region}"
  project = "slavayssiere-sandbox"
}

variable "region" {
  default = "europe-west1"
}

terraform {
  backend "gcs" {
    bucket = "tf-slavayssiere-wescale"
    prefix = "terraform/layer-istio-lb"
  }
}

resource "google_compute_global_address" "istio-lb-http" {
  name = "istio-lb-http"
}


data "google_dns_managed_zone" "public-gcp-wescale" {
  name = "slavayssiere-soa"
}


resource "google_dns_record_set" "istio-lb" {
  name = "iap.${data.google_dns_managed_zone.public-gcp-wescale.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = "${data.google_dns_managed_zone.public-gcp-wescale.name}"

  rrdatas = ["${google_compute_global_address.istio-lb-http.address}"]
}
