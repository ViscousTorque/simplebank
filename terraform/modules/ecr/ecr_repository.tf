module "ecr" {
  source = "terraform-aws-modules/ecr/aws"

  repository_name               = var.repository_name
  repository_image_scan_on_push = var.scan_on_push
  repository_lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep only the last 2 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 2
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
