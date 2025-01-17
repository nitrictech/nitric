
terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }

    azapi = {
      source = "Azure/azapi"
    }
  }
}

# data "docker_image" "latest" {
#   name = var.image_uri
# }

locals {
  remote_image_name = "${var.registry_login_server}/${var.stack_name}-${var.name}"
}

# Tag the provided docker image with the ECR repository url
resource "docker_tag" "tag" {
  source_image = var.image_uri
  target_image = local.remote_image_name
}

# Push the tagged image to the ECR repository
resource "docker_registry_image" "push" {
  name = local.remote_image_name
  triggers = {
    source_image_id = docker_tag.tag.source_image_id
  }
}


# Create an event token (random string) for the service
resource "random_string" "event_token" {
  length  = 32
  special = false
}

data "azuread_client_config" "current" {}

data "azurerm_client_config" "current" {}

locals {
  app_role_id    = "4962773b-9cdb-44cf-a8bf-237846a00ab7"
  repository_url = "${var.registry_login_server}/${var.stack_name}-${var.name}"
  # TODO: Remove these hardcoded role definitions
  role_definitions = {
    "KVSecretsOfficer"    = "b86a8fe4-44ce-4948-aee5-eccb2c155cd7"
    "BlobDataContrib"     = "ba92f5b4-2d11-453d-a403-e96b0029c9fe"
    "QueueDataContrib"    = "974c5e8b-45b9-4653-ba55-5f855dd0fb88"
    "EventGridDataSender" = "d5a91429-5739-47e2-a06b-3470a27159e7"
    "TagContributor"      = "4a9ae827-6dc8-4573-8ac7-8239d42aa03f"
  }
}

# Create a new azure ad service principal
resource "azuread_application" "service_identity" {
  display_name = "${var.name}-application-identity"
  app_role {
    allowed_member_types = ["Application"]
    description          = "Enables webhook subscriptions to authenticate using this application"
    display_name         = "AzureEventGrid Secure Webhook Subscriber"
    id                   = local.app_role_id
    value                = local.app_role_id
  }
  owners = [data.azuread_client_config.current.object_id]
}

resource "azuread_service_principal" "service_identity" {
  client_id = azuread_application.service_identity.client_id
  #   app_role_assignment_required = false
  owners = [data.azuread_client_config.current.object_id]
}

# Create a new app role assignment for the service principal
resource "azuread_app_role_assignment" "role_assignment" {
  app_role_id         = local.app_role_id
  principal_object_id = data.azuread_client_config.current.object_id
  resource_object_id  = azuread_service_principal.service_identity.object_id
}

# Create a new service principal password
resource "azuread_service_principal_password" "service_identity" {
  service_principal_id = azuread_service_principal.service_identity.id
}

# Assign roles to the service principal 
resource "azurerm_role_assignment" "role_assignment" {
  for_each = local.role_definitions

  principal_id       = azuread_service_principal.service_identity.id
  principal_type     = "ServicePrincipal"
  role_definition_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/${each.value}"
  scope              = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.resource_group_name}"
}

# Create a random string for the container app id
resource "random_string" "container_app_id" {
  length  = 4
  special = false
  upper   = false
}

locals {
  container_app_name = "${lower(replace(substr(var.name, 0, 28), "_", "-"))}-${random_string.container_app_id.result}"
}

# Create a new container app
# TODO...
resource "azurerm_container_app" "container_app" {
  name                         = local.container_app_name
  container_app_environment_id = var.container_app_environment_id
  resource_group_name          = var.resource_group_name
  revision_mode                = "Single"

  registry {
    server               = var.registry_login_server
    username             = var.registry_username
    password_secret_name = "registry-password"
  }

  ingress {
    external_enabled = true
    target_port      = 9001
    traffic_weight {
      latest_revision = true
      percentage = 100
    }
  }

  secret {
    name  = "registry-password"
    value = var.registry_password
  }

  secret {
    name  = "client-id"
    value = azuread_service_principal.service_identity.client_id
  }

  secret {
    name  = "client-secret"
    value = azuread_service_principal_password.service_identity.value
  }

  dapr {
    app_id       = local.container_app_name
    app_port     = 9001
    app_protocol = "http"
  }

  template {
    container {
      name   = "myapp"
      image  = docker_registry_image.push.name
      cpu    = var.cpu
      memory = var.memory

      env {
        name  = "EVENT_TOKEN"
        value = random_string.event_token.result
      }

      env {
        name  = "NITRIC_ENVIRONMENT"
        value = "cloud"
      }

      env {
        name  = "NITRIC_STACK_ID"
        value = var.stack_name
      }

      env {
        name  = "AZURE_SUBSCRIPTION_ID"
        value = data.azurerm_client_config.current.subscription_id
      }

      env {
        name  = "AZURE_RESOURCE_GROUP"
        value = var.resource_group_name
      }

      env {
        name  = "AZURE_CLIENT_ID"
        value = azuread_service_principal.service_identity.client_id
      }

      env {
        name        = "AZURE_CLIENT_SECRET"
        secret_name = "client-secret"
      }

      env {
        name  = "AZURE_TENANT_ID"
        value = data.azurerm_client_config.current.tenant_id
      }

      env {
        name  = "TOLERATE_MISSING_SERVICES"
        value = "true"
      }
      dynamic "env" {
        for_each = var.env
        content {
          name  = env.key
          value = env.value
        }
      }
    }
  }
}

resource "azapi_resource_action" "my_app_auth" {
  depends_on = [azurerm_container_app.container_app]

  type        = "Microsoft.App/containerApps/authConfigs@2024-03-01"
  resource_id = "${azurerm_container_app.container_app.id}/authConfigs/current"
  method      = "PUT"

  body = { # wrap in jsondecode if using 'azapi' v1
    location = azurerm_container_app.container_app.location
    properties = {
      globalValidation = {
        unauthenticatedClientAction = "Return401"
      }
      identityProviders = {
        azureActiveDirectory = {
          enabled = true
          registration = {
            clientId                =  azuread_application.service_identity.client_id
            clientSecretSettingName = "client-secret"
            openIdIssuer            = "https://sts.windows.net/${data.azuread_client_config.current.tenant_id}/v2.0"
          }
          validation = {
            allowedAudiences = azuread_application.service_identity.identifier_uris != null ? azuread_application.service_identity.identifier_uris : []
            # defaultAuthorizationPolicy = {
            #   allowedApplications = [
            #     azuread_application.my_app.client_id,
            #   ]
            # }
          }
        }
      }
      platform = {
        enabled = true
      }
    }
  }
}




