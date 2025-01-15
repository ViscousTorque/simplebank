variable "secret_name" {
  description = "Name of the Secrets Manager secret"
  type        = string
}

variable "description" {
  description = "Description of the secret"
  type        = string
}

variable "environment" {
  description = "Environment (e.g., dev, staging, prod)"
  type        = string
}

variable "db_source" {
  description = "Database source connection string"
  type        = string
}

variable "migration_url" {
  description = "Migration URL"
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
  description = "Enable reflection in gRPC services"
  type        = bool
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
  description = "Test recipient email address"
  type        = string
}

variable "email_sender_password" {
  description = "Email sender password"
  type        = string
}
