locals {
    neon_project_id = one(neon_project.project) != null ? one(neon_project.project).id : var.existing.project_id
    neon_branch_id = one(neon_branch.branch) != null ? one(neon_branch.branch).id : var.existing.branch_id
    neon_role_name = one(neon_role.role) != null ? one(neon_role.role).name : var.existing.role_name
    neon_role_password = data.neon_branch_role_password.password.password
    neon_endpoint_id = one(neon_endpoint.endpoint) != null ? one(neon_endpoint.endpoint).id : var.existing.endpoint_id
    neon_host_name = [for e in data.neon_branch_endpoints.endpoints.endpoints : e.host_name if e.id == neon_endpoint_id][0]
    neon_database_name = var.existing.database_name == null ? "${var.nitric.stack_id}-${var.nitric.name}" : var.existing.database_name
    neon_connection_string = "postgresql://${local.neon_role_name}:${local.neon_role_password}@${local.neon_host_name}/${local.neon_database_name}?sslmode=require"
}

resource "neon_project" "project" {
  count = var.existing.project_id == null ? 1 : 0
  name = "${var.nitric.stack_id}-${var.nitric.name}"
}

resource "neon_branch" "branch" {
  count = var.existing.branch_id == null ? 1 : 0
  project_id = local.neon_project_id
  name       = "${var.nitric.stack_id}-${var.nitric.name}"
}

resource "neon_role" "role" {
  count = var.existing.role_name == null ? 1 : 0

  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  name       = "${var.nitric.stack_id}-${var.nitric.name}"
}

resource "neon_database" "database" {
  count = var.existing.database_name == null ? 1 : 0

  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  name       = local.neon_database_name
  owner_name = local.neon_role_name
}

resource "neon_endpoint" "endpoint" {
  count = var.existing.endpoint_id == null ? 1 : 0

  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  type       = "read_write"
}

data "neon_branch_endpoints" "endpoints" {
  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
}

data "neon_branch_role_password" "password" {
  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  role_name  = local.neon_role_name
}