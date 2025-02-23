# Create a virtual network for the database server
resource "azurerm_virtual_network" "database_network" {
  name                = "nitric-database-vnet"
  resource_group_name = var.resource_group_name
  location            = var.location
  address_space       = ["10.0.0.0/16"]

  flow_timeout_in_minutes = 10
}

# Create a subnet for the database server
resource "azurerm_subnet" "database_subnet" {
  name                 = "nitric-database-subnet"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.database_network.name
  address_prefixes     = ["10.0.0.0/18"]

  delegation {
    name = "db-delegation"
    service_delegation {
      name    = "Microsoft.DBforPostgreSQL/flexibleServers"
    }
  }
}

# Create an infrastructure subnet for the database server
resource "azurerm_subnet" "database_infrastructure_subnet" {
  name                 = "nitric-database-infrastructure-subnet"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.database_network.name
  address_prefixes     = ["10.0.64.0/18"]

  depends_on = [ azurerm_subnet.database_subnet ]
}

# Create a subnet for containers to connect to the database
resource "azurerm_subnet" "database_client_subnet" {
  name                 = "nitric-database-client-subnet"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.database_network.name
  address_prefixes     = ["10.0.192.0/18"]

  delegation {
    name = "container-instance-delegation"
    service_delegation {
      name = "Microsoft.ContainerInstance/containerGroups"
    }
  }

  depends_on = [ azurerm_subnet.database_infrastructure_subnet ]
}

# Create a private zone for the database server
resource "azurerm_private_dns_zone" "database_dns_zone" {
  name                = "db-private-dns.postgres.database.azure.com"
  resource_group_name = var.resource_group_name
}

# Create a private link service for the database server
resource "azurerm_private_dns_zone_virtual_network_link" "database_link_service" {
  name                  = "nitric-database-link-service"
  private_dns_zone_name = azurerm_private_dns_zone.database_dns_zone.name
  resource_group_name   = var.resource_group_name
  virtual_network_id    = azurerm_virtual_network.database_network.id
  registration_enabled  = false
  tags = {
    "x-nitric-${var.stack_name}-name" = var.stack_name
    "x-nitric-${var.stack_name}-type" = "stack"
  }
}

# Create a random master password for the database server
resource "random_password" "database_master_password" {
  length  = 16
  special = false
}

# Create a database service if required
resource "azurerm_postgresql_flexible_server" "database" {
  name                         = "nitric-db-${random_string.stack_id.result}"
  resource_group_name          = var.resource_group_name
  location                     = var.location
  version                      = "14"
  administrator_login          = "nitric"
  administrator_password  = random_password.database_master_password.result

  public_network_access_enabled     = false
  
  delegated_subnet_id = azurerm_subnet.database_subnet.id
  private_dns_zone_id = azurerm_private_dns_zone.database_dns_zone.id

  # default to 32Gb storage
  # TODO: Make configurable   
  storage_mb = 32768

  # TODO: Make configurable  
  sku_name = "B_Standard_B1ms"

  tags = {
    "x-nitric-${var.stack_name}-name" = var.stack_name
    "x-nitric-${var.stack_name}-type" = "stack"
  }

  depends_on = [ 
    azurerm_subnet.database_subnet,
    azurerm_private_dns_zone.database_dns_zone,
    azurerm_private_dns_zone_virtual_network_link.database_link_service
  ]
}