variable "target_tags" {
  default = "tf-lb-https-gke"
}


variable "list_groups" {}

resource "google_compute_url_map" "my-url-map" {
  // note that this is the name of the load balancer
  name            = "my-url-map"
  default_service = "${module.gce-lb-https.backend_services[0]}"

  host_rule = {
    hosts        = ["*"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = "${module.gce-lb-https.backend_services[0]}"

    path_rule {
      paths   = ["/"]
      service = "${module.gce-lb-https.backend_services[0]}"
    }


    path_rule {
      paths   = ["/"]
      service = "${module.gce-lb-https.backend_services[1]}"
    }


    path_rule {
      paths   = ["/"]
      service = "${module.gce-lb-https.backend_services[2]}"
    }


    path_rule {
      paths   = ["/"]
      service = "${module.gce-lb-https.backend_services[3]}"
    }


    path_rule {
      paths   = ["/"]
      service = "${module.gce-lb-https.backend_services[4]}"
    }

    path_rule {
      paths   = ["/"]
      service = "${module.gce-lb-https.backend_services[5]}"
    }
  }

}

output "load-balancer-ip" {
  value = "${module.gce-lb-https.external_ip}"
}

module "gce-lb-https" {
  source            = "terraform-google-lb-http"
  name              = "istio-http"
  ssl               = false
  firewall_networks = ["demo-net"]

  // Make sure when you create the cluster that you provide the `--tags` argument to add the appropriate `target_tags` referenced in the http module. 
  target_tags = ["test-cluster"]

  // Use custom url map.
  url_map        = "${google_compute_url_map.my-url-map.self_link}"
  create_url_map = false

  // Get selfLink URLs for the actual instance groups (not the manager) of the existing GKE cluster:
  //   gcloud compute instance-groups list --uri
  backends = {
    "0" = [
      { group = "${element(split(";", var.list_groups),0)}" },
    ],
    "1" = [
      { group = "${element(split(";", var.list_groups),1)}" },
    ],
    "2" = [
      { group = "${element(split(";", var.list_groups),2)}" },
    ],
    "3" = [
      { group = "${element(split(";", var.list_groups),3)}" },
    ],
    "4" = [
      { group = "${element(split(";", var.list_groups),4)}" },
    ],
    "5" = [
      { group = "${element(split(";", var.list_groups),5)}" },
    ]
  }

  // You also must add the named port on the existing GKE clusters instance group that correspond to the `service_port` and `service_port_name` referenced in the module definition.
  //   gcloud compute instance-groups set-named-ports INSTANCE_GROUP_NAME --named-ports=NAME:PORT
  // replace `INSTANCE_GROUP_NAME` with the name of your GKE cluster's instance group and `NAME` and `PORT` with the values of `service_port_name` and `service_port` respectively.
  backend_params = [
    "/,http,31380,10",
    "/,http,31380,10",
    "/,http,31380,10",
    "/,http,31380,10",
    "/,http,31380,10",
    "/,http,31380,10",
  ]
}
