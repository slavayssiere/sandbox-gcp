resource "google_container_cluster" "test-cluster" {
  provider           = "google-beta"
  name               = "test-cluster"
  region             = "${var.region}"
  initial_node_count = 1

  private_cluster_config {
    enable_private_endpoint = true
    enable_private_nodes    = true
    master_ipv4_cidr_block = "192.168.16.0/28"
  }

  min_master_version = "1.11.2-gke.18"
  node_version       = "1.11.2-gke.18"

  network    = "demo-net"
  subnetwork = "demo-subnet"

  addons_config {
    kubernetes_dashboard {
      disabled = true
    }
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = "c0-pods"
    services_secondary_range_name = "c0-services"
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels {
      Name = "test-cluster"
    }

    tags = ["kubernetes", "test-cluster"]
  }
}
