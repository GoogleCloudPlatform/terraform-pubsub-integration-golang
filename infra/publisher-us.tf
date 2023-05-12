locals {
  us_west1_publisher_namespace                = "us-west1-publisher"
  us_west1_publisher_k8s_service_account_name = "us-west1-publisher"
  us_west1_publisher_base_entries = [
    {
      name  = "namespace"
      value = local.us_west1_publisher_namespace
    },
    {
      name  = "gcp_service_account_email"
      value = module.us_west1_publisher_cluster.gcp_service_account_email
    },
    {
      name  = "k8s_service_account_name"
      value = local.us_west1_publisher_k8s_service_account_name
    },
  ]
}

module "us_west1_publisher_cluster" {
  depends_on = [
    module.project_services,
  ]
  source = "./modules/kubernetes"

  cluster_name           = "us-west1-publisher-golang"
  region                 = "us-west1"
  zones                  = ["us-west1-a"]
  network_self_link      = google_compute_network.primary.self_link
  project_id             = data.google_project.project.project_id
  gcp_service_account_id = "us-west1-publisher-golang"
  gcp_service_account_iam_roles = [
    "roles/pubsub.publisher",
  ]
  k8s_namespace_name       = local.us_west1_publisher_namespace
  k8s_service_account_name = local.us_west1_publisher_k8s_service_account_name
  labels                   = var.labels
}

module "us_west1_publisher_base_helm" {
  source = "./modules/helm"

  providers = {
    helm = helm.us_west1_publisher_helm
  }
  chart_folder_name = "base"
  region            = "us-west1"
  entries           = local.us_west1_publisher_base_entries
}

module "us_west1_publisher_helm" {
  depends_on = [
    module.us_west1_publisher_base_helm,
  ]
  source = "./modules/helm"

  providers = {
    helm = helm.us_west1_publisher_helm
  }
  chart_folder_name = "publisher"
  region            = "us-west1"
  entries = concat(local.us_west1_publisher_base_entries,
    [
      {
        name  = "project_id"
        value = data.google_project.project.project_id
      },
      {
        name  = "region"
        value = "us-west1"
      },
      {
        name  = "image"
        value = var.publisher_image_url
      },
      {
        name  = "config_maps.event_topic"
        value = google_pubsub_topic.event.id
      },
    ]
  )
}
