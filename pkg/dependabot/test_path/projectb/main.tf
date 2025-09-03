terraform {
  backend "gcs" {
    bucket = "terraform-state-iac-cloudflare"
    prefix = "organization"
  }
}
