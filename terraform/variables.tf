variable "description" {
  description = "Description of the secret"
  type        = string
}

variable "secret_name" {
  description = "Name of the secret"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "db_source" {
  description = "Database connection string"
  type        = string
}

variable "migration_url" {
  description = "Database migration URL"
  type        = string
}

variable "http_server_address" {
  description = "HTTP server address"
  type        = string
}

variable "grpc_server_address" {
  description = "gRPC server address"
  type        = string
}

variable "redis_address" {
  description = "Redis server address"
  type        = string
}

variable "token_symmetric_key" {
  description = "Token symmetric key"
  type        = string
}

variable "access_token_duration" {
  description = "Access token duration"
  type        = string
}

variable "refresh_token_duration" {
  description = "Refresh token duration"
  type        = string
}

variable "enable_reflection" {
  description = "Enable reflection (true/false)"
  type        = string
}

variable "email_sender_name" {
  description = "Email sender name"
  type        = string
}

variable "email_sender_address" {
  description = "Email sender address"
  type        = string
}

variable "email_test_recipient" {
  description = "Email test recipient address"
  type        = string
}

variable "email_sender_password" {
  description = "Email sender password"
  type        = string
  sensitive   = true
}

variable "repository_name" {
  description = "Name of the ECR repository"
  type        = string
}

variable "image_tag_mutability" {
  description = "The tag mutability setting for the repository (MUTABLE or IMMUTABLE)"
  type        = string
}

variable "scan_on_push" {
  description = "Indicates whether images are scanned on push"
  type        = bool
}

variable "encryption_type" {
  description = "The encryption type to use for the repository (AES256 or KMS)"
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC where the security group will be created"
  type        = string
}

variable "sg_ingress_cidr" {
  description = "CIDR blocks allowed for ingress"
  type        = string
}

variable "project_name" {
  description = "Name for this project"
  type        = string
}

variable "rds_db_user_password" {
  description = "Non IAM authenticated db user password"
  type        = string
}


