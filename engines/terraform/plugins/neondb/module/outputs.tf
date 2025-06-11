output "nitric" {
    value = {
        id          = neon_endpoint.endpoint.id
        exports = {
            # Export known service outputs
            services = local.service_outputs
            resources = {}
        }
    } 
}