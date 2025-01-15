variable "rds_db_username" {
  description = "db username"
  type        = string
}

variable "rds_db_user_password" {
  description = "Non IAM authenticated db user password"
  type        = string
}

variable "vpc_security_group_ids" {
  description = "The list of security group ids set on the rds"
  type        = list(string)
}

variable "rds_subnet_ids" {
  description = "The list of subnet ids set on the rds"
  type        = list(string)
}

variable "rds_db_subnet_group_name" {
  description = "Name of DB subnet group. DB instance will be created in the VPC associated with the DB subnet group"
  type        = string
}

variable "rds_availability_zone" {
  description = "The availablity zone the free tier rds is operating in"
  type        = string
}

variable "db_name" {
  description = "The sql db name in the postgres instance"
  type        = string
}

variable "engine" {
  description = "The engine type for the rds instance"
  type        = string
}

variable "rds_db_family" {
  description = "The family of sql DB parameter group"
  type        = string
}

variable "engine_version" {
  description = "The engine version for the rds instance"
  type        = string
}

variable "instance_identifier" {
  description = "The RDS DB instance name"
  type        = string
}      