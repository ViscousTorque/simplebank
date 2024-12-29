module "simple_bank_secret" {
  source  = "terraform-aws-modules/secrets-manager/aws"
  version = "1.3.1"

  name                    = var.secret_name
  description             = var.description
  recovery_window_in_days = 0

  secret_string = jsonencode({
    ENVIRONMENT            = var.environment
    DB_SOURCE              = var.db_source
    MIGRATION_URL          = var.migration_url
    HTTP_SERVER_ADDRESS    = var.http_server_address
    GRPC_SERVER_ADDRESS    = var.grpc_server_address
    REDIS_ADDRESS          = var.redis_address
    TOKEN_SYMMETRIC_KEY    = var.token_symmetric_key
    ACCESS_TOKEN_DURATION  = var.access_token_duration
    REFRESH_TOKEN_DURATION = var.refresh_token_duration
    ENABLE_REFLECTION      = var.enable_reflection
    EMAIL_SENDER_NAME      = var.email_sender_name
    EMAIL_SENDER_ADDRESS   = var.email_sender_address
    EMAIL_TEST_RECIPIENT   = var.email_test_recipient
    EMAIL_SENDER_PASSWORD  = var.email_sender_password
  })

  tags = {
    Environment = var.environment
    Project     = "SimpleBank"
  }
}



