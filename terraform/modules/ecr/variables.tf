variable "repository_name" {
  description = "Name of the ECR repository"
  type        = string
  default     = "simple_bank"
}

variable "image_tag_mutability" {
  description = "The tag mutability setting for the repository (MUTABLE or IMMUTABLE)"
  type        = string
  default     = "MUTABLE"
}

variable "scan_on_push" {
  description = "Indicates whether images are scanned on push"
  type        = bool
  default     = false
}

variable "encryption_type" {
  description = "The encryption type to use for the repository (AES256 or KMS)"
  type        = string
  default     = "AES256"
}

