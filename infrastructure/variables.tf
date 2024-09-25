variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The region to deploy to"
  type        = string
  default     = "us-central1"
}

variable "credentials" {
  description = "The GCP service account JSON credentials"
  type        = string
  sensitive   = true
}
