provider "aws" {
  region = "us-east-1"
}

module "ec2-east" {
   source = "./ec2"
   name = "ec2-htmx-agario"
   ami = "ami-07d9b9ddc6cd8dd30"
}





