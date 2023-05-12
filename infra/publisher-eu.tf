locals {
  europe_north1_publisher_namespace                = "europe-north1-publisher"
  europe_north1_publisher_k8s_service_account_name = "europe-north1-publisher"
  europe_north1_base_entries = [
    {
      name  = "namespace"
      value = local.europe_north1_publisher_namespace
    },
    {
      name  = "gcp_service_account_email"
      value = module.europe_north1_publisher_cluster.gcp_service_account_email
    },
    {
      name  = "k8s_service_account_name"
      value = local.europe_north1_publisher_k8s_service_account_name
    },
  ]
}

module "europe_north1_publisher_cluster" {
  depends_on = [
    module.project_services,
  ]
  source = "./modules/kubernetes"

  cluster_name           = "europe-north1-publisher-golang"
  region                 = "europe-north1"
  zones                  = ["europe-north1-a"]
  network_self_link      = google_compute_network.primary.self_link
  project_id             = data.google_project.project.project_id
  gcp_service_account_id = "europe-north1-publisher-golang"
  gcp_service_account_iam_roles = [
    "roles/pubsub.publisher",
  ]
  k8s_namespace_name       = local.europe_north1_publisher_namespace
  k8s_service_account_name = local.europe_north1_publisher_k8s_service_account_name
  labels                   = var.labels
}


module "europe_north1_publisher_base_helm" {
  source = "./modules/helm"

  providers = {
    helm = helm.europe_north1_publisher_helm
  }
  chart_folder_name = "base"
  region            = "europe-north1"
  entries           = local.europe_north1_base_entries
}

module "europe_north1_publisher_helm" {
  depends_on = [
    module.europe_north1_publisher_base_helm,
  ]
  source = "./modules/helm"

  providers = {
    helm = helm.europe_north1_publisher_helm
  }
  chart_folder_name = "publisher"
  region            = "europe-north1"
  entries = concat(local.europe_north1_base_entries,
    [
      {
        name  = "project_id"
        value = data.google_project.project.project_id
      },
      {
        name  = "region"
        value = "europe-north1"
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
