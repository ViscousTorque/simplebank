
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = var.project_name
  cidr = var.vpc_cidr

  azs = var.azs

  public_subnets  = [for k, v in var.azs : cidrsubnet(var.vpc_cidr, 8, k)]
  private_subnets = [for k, v in var.azs : cidrsubnet(var.vpc_cidr, 8, k + 3)]

  enable_dns_support   = true
  enable_dns_hostnames = true
  enable_nat_gateway   = true
  single_nat_gateway   = true
}

// TODO - this resource is against best practices, only good for quick demo, fix later for best practices
resource "aws_db_subnet_group" "public" {
  name        = "public-db-subnet-group"
  description = "DB Subnet Group for Public Subnets"
  subnet_ids  = module.vpc.public_subnets

  tags = {
    Name = "Public DB Subnet Group"
  }
}


