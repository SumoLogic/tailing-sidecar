#################################################################################
# Variables
#################################################################################

variable "TAGS" {
  type = set(string)
  default = [
    "local/tailing-sidecar-otel:main",
    "local/tailing-sidecar-otel:${VERSION}",
  ]

  validation {
    condition = TAGS != ""
    error_message = "The variable 'TAGS' must not be empty."
  }
}

# NOTE: Must be passed the output of the command: git describe --tags --always
# Example value: v0.17.2
variable "VERSION" {
  type = string

  validation {
    condition = VERSION != ""
    error_message = "The variable 'VERSION' must not be empty."
  }
}

#################################################################################
# Groups
#################################################################################

group "default" {
  targets = [
    "local"
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
    VERSION = "${VERSION}"
  }
  dockerfile = "./Dockerfile"
}

target "_common-local" {
  inherits = ["_common"]
  output = [{ type = "docker" }]
}

target "_common-multiplatform" {
  inherits = ["_common"]
  platforms = [
    "linux/amd64",
    "linux/arm64"
  ]
}

#################################################################################
# Composite targets
#################################################################################

target "local" {
  inherits = ["_common-local"]
  tags = [
    "local/tailing-sidecar-otel:main",
    "local/tailing-sidecar-otel:${VERSION}"
  ]
}

target "test" {
  inherits = ["_common-multiplatform"]
  tags = [
    "ghcr.io/sumologic/tailing-sidecar-otel:main",
    "ghcr.io/sumologic/tailing-sidecar-otel:${VERSION}"
  ]
}

target "dev" {
  inherits = ["_common-multiplatform"]
  tags = [
    "public.ecr.aws/sumologic/tailing-sidecar-otel-dev:main",
    "public.ecr.aws/sumologic/tailing-sidecar-otel-dev:${VERSION}",
    "ghcr.io/sumologic/tailing-sidecar-otel:main",
    "ghcr.io/sumologic/tailing-sidecar-otel:${VERSION}",
    "sumologic/tailing-sidecar-otel-dev:main",
    "sumologic/tailing-sidecar-otel-dev:${VERSION}"
  ]
}

target "prod" {
  inherits = ["_common-multiplatform"]
  tags = [
    "public.ecr.aws/sumologic/tailing-sidecar-otel-dev:main",
    "public.ecr.aws/sumologic/tailing-sidecar-otel-dev:${VERSION}",
    "ghcr.io/sumologic/tailing-sidecar-otel:main",
    "ghcr.io/sumologic/tailing-sidecar-otel:${VERSION}",
    "sumologic/tailing-sidecar-otel-dev:main",
    "sumologic/tailing-sidecar-otel-dev:${VERSION}"
  ]
}
