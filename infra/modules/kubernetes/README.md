# KUBERNETES

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cluster\_name | The name of the cluster | `string` | n/a | yes |
| gcp\_service\_account\_iam\_roles | the list of permissions for gcp service account | `list(string)` | n/a | yes |
| gcp\_service\_account\_id | gcp service account's id | `string` | n/a | yes |
| k8s\_namespace\_name | kubernetes namespace name | `string` | n/a | yes |
| k8s\_service\_account\_name | kubernetes service account for pods | `string` | n/a | yes |
| labels | A map of key/value label pairs to assign to the resources. | `map(string)` | `{}` | no |
| network\_self\_link | The VPC network self\_link to host the k8s cluster | `string` | n/a | yes |
| project\_id | GCP project ID. | `string` | n/a | yes |
| region | The region to host the k8s cluster | `string` | n/a | yes |
| zones | The zones to host the k8s cluster | `list(string)` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| cluster\_info | GKE cluster information containing name, location and namespace |
| control\_plane | Control plane in deployed GKE cluster |
| gcp\_service\_account\_email | Service account’s email address in the project |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
