
module "security_group" {
  source      = "terraform-aws-modules/security-group/aws"
  version     = "5.2.0"
  name        = var.sg_name
  description = var.sg_description
  vpc_id      = var.vpc_id

  ingress_with_cidr_blocks = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      cidr_blocks = var.sg_ingress_cidr
    }
  ]

  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

}

output "security_group_id" {
  value = module.security_group.security_group_id
}
