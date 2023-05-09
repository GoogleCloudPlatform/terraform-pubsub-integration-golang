terraform {
  required_version = ">= 0.13"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.53"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.9"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "helm" {
  alias = "europe_north1_publisher_helm"
  kubernetes {
    host                   = "https://${module.europe_north1_publisher_cluster.control_plane.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(module.europe_north1_publisher_cluster.control_plane.master_auth[0].cluster_ca_certificate, )
    client_certificate     = base64decode(module.europe_north1_publisher_cluster.control_plane.master_auth[0].client_certificate)
    client_key             = base64decode(module.europe_north1_publisher_cluster.control_plane.master_auth[0].client_key)
  }
}

provider "helm" {
  alias = "us_west1_publisher_helm"
  kubernetes {
    host                   = "https://${module.us_west1_publisher_cluster.control_plane.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(module.us_west1_publisher_cluster.control_plane.master_auth[0].cluster_ca_certificate, )
    client_certificate     = base64decode(module.us_west1_publisher_cluster.control_plane.master_auth[0].client_certificate)
    client_key             = base64decode(module.us_west1_publisher_cluster.control_plane.master_auth[0].client_key)
  }
}

provider "helm" {
  alias = "us_west1_subscriber_helm"
  kubernetes {
    host                   = "https://${module.us_west1_subscriber_cluster.control_plane.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(module.us_west1_subscriber_cluster.control_plane.master_auth[0].cluster_ca_certificate, )
    client_certificate     = base64decode(module.us_west1_subscriber_cluster.control_plane.master_auth[0].client_certificate)
    client_key             = base64decode(module.us_west1_subscriber_cluster.control_plane.master_auth[0].client_key)
  }
}
