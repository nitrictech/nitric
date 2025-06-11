output "nitric" {
    value = {
        id          = local.neon_endpoint_id
        exports = {
            # Export known service outputs
            services = local.service_outputs
            resources = {}
        }
    } 
}