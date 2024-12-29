output "vpc_id" {
  value = module.vpc.vpc_id
}

output "db_subnet_group_name" {
  value       = aws_db_subnet_group.public.name
  description = "Hacking the public subnet to use for easy access rds"
}

output "db_subnet_ids" {
  value       = aws_db_subnet_group.public.subnet_ids
  description = "Hacking the public subnet to use for easy access rds"
}