output "image_id" {
  description = "A reference to the locally built image"
  value = docker_image.service.image_id
}
