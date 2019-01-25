resource "google_storage_bucket" "code-dataflow-bucket" {
  name     = "code-dataflow-bucket"
  location = "europe-west1"
  storage_class = "REGIONAL"
  force_destroy = true
}