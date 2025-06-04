output "nitric" {
  value = {
    role = aws_iam_role.role
    id   = aws_iam_role.role.name
  }
}
