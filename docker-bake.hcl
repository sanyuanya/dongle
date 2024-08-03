variable "TAG" {
  default = "v0.0.1"
}

group "default" {
  targets = ["dongle"]
}

target "dongle" {
  dockerfile = "Dockerfile"
  context = "."
  tags = ["hj212223/dongle:${TAG}"]
  platforms = ["linux/amd64", "linux/arm64"]
}