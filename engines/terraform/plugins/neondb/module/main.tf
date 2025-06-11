locals {
    neon_project_id = one(neon_project.project) != null ? one(neon_project.project).id : var.existing.project_id
    neon_branch_id = one(neon_branch.branch) != null ? one(neon_branch.branch).id : var.existing.branch_id
    # neon_role_name = one(neon_role.role) != null ? one(neon_role.role).name : var.existing.role_name
    # neon_role_password = data.neon_branch_role_password.password.password
    # neon_endpoint_id = one(neon_endpoint.endpoint) != null ? one(neon_endpoint.endpoint).id : null
    # neon_host_name =  [for e in data.neon_branch_endpoints.endpoints.endpoints : e.host_name if e.id == local.neon_endpoint_id][0]
    # neon_database_name = var.existing.database_name == null ? "${var.nitric.stack_id}-${var.nitric.name}" : var.existing.database_name
    neon_database_name = "${var.nitric.stack_id}-${var.nitric.name}"
    neon_connection_string = "postgresql://${neon_role.role.name}:${neon_role.role.password}@${neon_endpoint.endpoint.host}/${local.neon_database_name}?sslmode=require"

    # Output service export map
    service_outputs = {
        for name, service in var.nitric.services : name => {
            env = {
                DATABASE_URL = local.neon_connection_string
            }
        }
    }
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


resource "neon_endpoint" "endpoint" {
  # count = var.existing.endpoint_id == null ? 1 : 0

  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  type       = "read_write"
}
resource "neon_role" "role" {
  # count = var.existing.role_name == null ? 1 : 0
  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  name       = "${var.nitric.stack_id}-${var.nitric.name}"

  depends_on = [ neon_endpoint.endpoint ]
}

resource "neon_database" "database" {
  # count = var.existing.database_name == null ? 1 : 0

  project_id = local.neon_project_id
  branch_id  = local.neon_branch_id
  name       = local.neon_database_name
  owner_name = neon_role.role.name
}



# data "neon_branch_endpoints" "endpoints" {
#   project_id = local.neon_project_id
#   branch_id  = local.neon_branch_id
# }

# data "neon_branch_role_password" "password" {
#   project_id = local.neon_project_id
#   branch_id  = local.neon_branch_id
#   role_name  = local.neon_role_name
# }