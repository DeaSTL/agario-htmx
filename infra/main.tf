provider "aws" {
  region = "us-east-1"
}

module "ec2-east" {
   source = "./ec2"
   name = "htmx-agario"
   ami = "ami-07d9b9ddc6cd8dd30"
   route_zone_id = "Z05859407AHEY4L0BH7B"
   domain = "aws.jmhart.dev"
}





