resource "docker_image" "base_service" {
  name = var.image_id == null ? "base_service" : var.image_id

  dynamic "build" {
    for_each = var.image_id == null ? [1] : []
    content {
      # TODO: 
      context = "${path.root}/../../${var.build_context != "." ? var.build_context : ""}"
      dockerfile = var.dockerfile
      tag     = ["${var.tag}:base"]
    }
  }
}

# Extract entrypoint and command using Docker CLI via external data source
data "external" "inspect_base_image" {
  depends_on = [docker_image.base_service]
  program = ["docker", "inspect", docker_image.base_service.image_id, "--format", "{\"entrypoint\":\"{{join .Config.Entrypoint \" \"}}\",\"cmd\":\"{{join .Config.Cmd \" \"}}\"}" ]
}

locals {
  original_command = "${data.external.inspect_base_image.result.entrypoint} ${data.external.inspect_base_image.result.cmd}"
  image_id = var.image_id == null ? docker_image.base_service.name : var.image_id
}

# Next we want to wrap this image withing a nitric service
resource "docker_image" "service" {
  name = "service"
  build {
    # This doesn't actually matter as we aren't copying in anything relative
    context = "."
    # Use the wrapped dockerfile here
    dockerfile = "${path.module}/wrapped.dockerfile"
    build_args = merge({
      BASE_IMAGE = local.image_id
      ORIGINAL_COMMAND = local.original_command
    }, var.args)
    tag     = ["${var.tag}:latest"]
  }
}