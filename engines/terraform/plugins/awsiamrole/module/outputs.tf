output "nitric" {
  value = {
    # FIXME: Need to make output type consistent with other identity plugins
    # to avoid errors when mixing them as inputs to other modules
    role = aws_iam_role.role
    id   = aws_iam_role.role.name
  }
}
