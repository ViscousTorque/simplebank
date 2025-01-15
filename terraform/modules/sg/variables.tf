variable "vpc_id" {
  description = "The ID of the VPC where the security group will be created"
  type        = string
}

variable "sg_ingress_cidr" {
  description = "CIDR blocks allowed for ingress"
  type        = string
}

variable "sg_description" {
  description = "A description for the security group"
  type        = string
  default     = "Security managed with Terraform"
}

variable "sg_name" {
  description = "The name of the security group"
  type        = string
}
