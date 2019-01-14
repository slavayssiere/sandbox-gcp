resource "google_bigtable_instance" "test-instance" {
  name         = "test-instance"
  cluster {
    cluster_id   = "tf-instance-cluster"
    zone         = "europe-west1-b"
    num_nodes    = 1
    storage_type = "HDD"
  }
}

resource "google_bigtable_table" "test-table" {
  name          = "test-table"
  instance_name = "${google_bigtable_instance.test-instance.name}"
  # split_keys    = ["a", "b", "c"]
}

// gcloud beta bigtable instances create $BT_INSTANCE \
//   --cluster=$BT_INSTANCE \
//   --cluster-zone="europe-west1-b" \
//   --display-name=$BT_INSTANCE \
//   --cluster-num-nodes=3

// # Sample code to bootstrap BigTable structure and inserting a sample value
// cbt -project $PROJECT -instance $BT_INSTANCE createtable my-table

// cbt -instance $BT_INSTANCE createfamily my-table tests
