#################################################################################
# NOTE: This bake file was added to make it easier to test variations of the
# docker images locally without having to update the Makefile. It isn't used in
# any part of the release process.
#################################################################################

#################################################################################
# Variables
#################################################################################

variable "RELEASE_NUMBER" {
  type = string
  default = "1"
}

variable "TAG" {
  type = string
  default = "localhost:32000/sumologic/tailing-sidecar:latest"
}

variable "VERSION" {
  type = string
  default = ""
}

#################################################################################
# Groups
#################################################################################

group "default" {
  targets = [
    "standard",
    "ubi"
  ]
}

#################################################################################
# Base targets
#################################################################################

target "_common" {
  context = "./"
  attest = [
    {
      type = "provenance",
      disabled = true,
    },
    {
      type = "sbom",
      disabled = true,
    },
  ]
  args = {
    RELEASE_NUMBER = "${RELEASE_NUMBER}"
    VERSION = "${VERSION}"
  }
  output = [{ type = "docker" }]
}

target "_common-standard" {
  inherits = ["_common"]
  dockerfile = "./Dockerfile"
}

target "_common-ubi" {
  inherits = ["_common"]
  dockerfile = "./Dockerfile.ubi"
}

#################################################################################
# Composite targets
#################################################################################

target "standard" {
  inherits = ["_common-standard"]
  tags = ["${TAG}"]
}

# NOTE: Only linux/amd64 is supported as our fluent-bit base image was not
# pushed with support for other platforms. This may be something we can change
# but it may not make sense as we're moving away from fluent-bit.
target "ubi" {
  inherits = ["_common-ubi"]
  tags = ["${TAG}-ubi"]
  platforms = ["linux/amd64"]
}
