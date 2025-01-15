
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.80.0"
    }
  }
}

module "secretsmanager" {
  source = "./modules/secretsmanager"

  description            = "Simple bank secret"
  secret_name            = "simple_bank"
  environment            = var.environment
  db_source              = var.db_source
  migration_url          = var.migration_url
  http_server_address    = var.http_server_address
  grpc_server_address    = var.grpc_server_address
  redis_address          = var.redis_address
  token_symmetric_key    = var.token_symmetric_key
  access_token_duration  = var.access_token_duration
  refresh_token_duration = var.refresh_token_duration
  enable_reflection      = var.enable_reflection
  email_sender_name      = var.email_sender_name
  email_sender_address   = var.email_sender_address
  email_test_recipient   = var.email_test_recipient
  email_sender_password  = var.email_sender_password
}

module "ecr" {
  source               = "./modules/ecr"
  repository_name      = var.repository_name
  image_tag_mutability = var.image_tag_mutability
  scan_on_push         = var.scan_on_push
  encryption_type      = var.encryption_type
}

module "security_group" {
  source          = "./modules/sg"
  sg_name         = "access-postgres-anywhere"
  sg_description  = "RDS postgres security group for demo purposes"
  sg_ingress_cidr = var.sg_ingress_cidr
  vpc_id          = module.vpc.vpc_id
}

module "vpc" {
  source   = "./modules/vpc"
  azs      = local.azs
  vpc_cidr = local.vpc_cidr
  project_name = "simplebank"
}

module "rds" {
  source                   = "./modules/rds"
  rds_availability_zone    = local.azs[0]
  db_name                  = "simplebank"
  engine                   = "postgres"
  rds_db_family            = "postgres16"
  engine_version           = "16.3"
  instance_identifier      = "simple-bank-db-inst"
  rds_db_username          = "postgresAdmin"
  rds_db_user_password     = var.rds_db_user_password
  rds_subnet_ids           = module.vpc.db_subnet_ids
  rds_db_subnet_group_name = module.vpc.db_subnet_group_name
  vpc_security_group_ids   = [module.security_group.security_group_id]
}


