resource "google_service_account" "bastion-sa" {
  account_id   = "bastion-sa"
  display_name = "Bastion Service Account"
}

// resource "google_project_iam_custom_role" "admin-rbac" {
//   role_id     = "kubeadmin"
//   title       = "Admin RBAC"
//   description = "For Kubernetes bastion"
//   permissions = ["container.clusterRoleBindings.create",
//                 "container.clusterRoleBindings.delete",
//                 "container.clusterRoleBindings.get",
//                 "container.clusterRoleBindings.list",
//                 "container.clusterRoleBindings.update",
//                 "container.clusters.get",
//                 "container.clusters.list",
//                 "container.clusters.getCredentials",
//                 "container.clusters.create",
//                 "container.clusters.delete",
//                 "container.clusters.update",
//                 "container.operations.get",
//                 "container.operations.list"
//                 ]
// }

// data "google_iam_policy" "kubernetes-admin" {
//   binding {
//     role = "projects/slavayssiere-sandbox/roles/kubeadmin"
//     members = ["serviceAccount:bastion-sa@slavayssiere-sandbox.iam.gserviceaccount.com"]
//   }
// }

// resource "google_service_account_iam_policy" "bastion-sa-iam" {
//     service_account_id = "${google_service_account.bastion-sa.name}"
//     policy_data = "${data.google_iam_policy.kubernetes-admin.policy_data}"
// }

resource "google_service_account_iam_binding" "admin-account-iam" {
  service_account_id = "${google_service_account.bastion-sa.name}"
  role        = "projects/slavayssiere-sandbox/roles/admin_kubernetes_bastion"

  members = [
    "serviceAccount:bastion-sa@slavayssiere-sandbox.iam.gserviceaccount.com",
  ]
}