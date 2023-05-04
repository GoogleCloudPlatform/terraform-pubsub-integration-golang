output "project_id" {
  description = "GCP project ID."
  value       = data.google_project.current.project_id
}

output "errors_topic_name" {
  description = "The name of the error topic"
  value       = google_pubsub_topic.errors.name
}

output "metrics_topic_name" {
  description = "The name of the metric topic"
  value       = google_pubsub_topic.metrics.name
}

output "event_topic_name" {
  description = "The name of the event topic"
  value       = google_pubsub_topic.event.name
}
