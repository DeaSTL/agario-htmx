
resource "aws_vpc" "htmx-agario-vpc" {
    cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "htmx-agario-subnet" {
    vpc_id            = aws_vpc.htmx-agario-vpc.id
    availability_zone = "us-east-1a"
    cidr_block = "10.0.0.0/16"
    map_public_ip_on_launch = true
}

resource "aws_security_group" "htmx-agario-sg" {
  name = "htmx-agario-sg"
  description = "Allow 8080,443,80"

  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    name = "htmx-agario-sg"
  }
}

resource "aws_key_pair" "htmx_agario" {
  key_name = "htmx-agario-key"
  public_key = file("~/.ssh/id_rsa.pub")
}

resource "aws_instance" "htmx-agario-instance" {
  ami = var.ami
  instance_type = "t2.micro"
  key_name = aws_key_pair.htmx_agario.key_name
  tags = {
    Name = var.name
  }
  security_groups = [aws_security_group.htmx-agario-sg.name]

  user_data = templatefile("${path.module}/cloud-init.yaml",{ssh_key:"${aws_key_pair.htmx_agario.public_key}"})

}
