terraform {
  backend "s3" {
    bucket = "alberto-secure-pipeline-tf-state"
    key    = "af-terraform-state"
    region = "eu-west-2"
  }
}