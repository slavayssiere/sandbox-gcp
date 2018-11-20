
provider "google" {
  region      = "${var.region}"
  project = "slavayssiere-sandbox"
}


provider "google-beta" {
  region      = "${var.region}"
  project = "slavayssiere-sandbox"
}

variable "region" {
    default = "europe-west1"
}


terraform {
  backend "gcs" {
    bucket  = "tf-slavayssiere-wescale"
    prefix  = "terraform/layer-kubernetes"
    region = "${var.region}"
  }
}

data "terraform_remote_state" "layer-base" {
  backend = "gcp"
  config {
    bucket  = "tf-slavayssiere-wescale"
    prefix  = "terraform/layer-base"
  }
}



