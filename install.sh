gcloud builds submit --tag gcr.io/roi-takeoff-user3/events-website:v1.0.5
terraform init && terraform apply -auto-approve --var="project_id=roi-takeoff-user3"
