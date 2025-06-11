output "nitric" {
    value = {
        id          = local.neon_endpoint_id
        exports = {
            env = {
                "DB_URL" = local.neon_connection_string
            }
            resources = {}
        }
    } 
}