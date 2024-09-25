# Google Cloud Provider Configuration
provider "google" {
  project     = var.project_id
  region      = var.region
  credentials = var.credentials
}

# Create a Google Cloud Storage bucket
resource "google_storage_bucket" "airportima_bucket" {
  name          = "airportima-bucket"
  location      = var.region
  force_destroy = true
  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  labels = {
    environment = "production"
    purpose     = "airport-image-storage"
  }
}

# Outputs to provide information about the created resources
output "bucket_name" {
  description = "The name of the created GCS bucket"
  value       = google_storage_bucket.airportima_bucket.name
}

output "bucket_url" {
  description = "The URL of the created GCS bucket"
  value       = "gs://${google_storage_bucket.airportima_bucket.name}"
}
