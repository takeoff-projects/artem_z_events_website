provider "google" {
  project = var.project_id
  region  = var.provider_region
}

resource "google_cloud_run_service" "default" {
  name     = "cloudrun-events"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/roi-takeoff-user3/events-website:v1.0.2"
        env {
          name = "GOOGLE_CLOUD_PROJECT"
          value = var.project_id
      }
    }
  }
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location    = google_cloud_run_service.default.location
  project     = google_cloud_run_service.default.project
  service     = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}