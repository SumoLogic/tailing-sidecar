#!/usr/bin/env bash

#################################################################################
# This script reads a builder config template, replaces the dist version and
# otel versions, and then outputs a builder config.
#################################################################################

set -euo pipefail

#################################################################################
# Command options
#################################################################################

OPT_SHORT_HELP="h"
OPT_LONG_HELP="help"
readonly OPT_SHORT_HELP OPT_LONG_HELP

OPT_SHORT_INPUT="i"
OPT_LONG_INPUT="input"
readonly OPT_SHORT_INPUT OPT_LONG_INPUT

OPT_SHORT_OUTPUT="o"
OPT_LONG_OUTPUT="output"
readonly OPT_SHORT_OUTPUT OPT_LONG_OUTPUT

#################################################################################
# Default values
#################################################################################

DEFAULT_INPUT_PATH="builder-template.yaml"
DEFAULT_OUTPUT_PATH="builder-config.yaml"
DEFAULT_DIST_VERSION=""
if command -v git; then
    DEFAULT_DIST_VERSION="$(git describe --tags --always)"
fi

INPUT_PATH="${INPUT_PATH:-$DEFAULT_INPUT_PATH}"
OUTPUT_PATH="${OUTPUT_PATH:-$DEFAULT_OUTPUT_PATH}"
DIST_VERSION="${DIST_VERSION:-$DEFAULT_DIST_VERSION}"
OT_VERSION="${OT_VERSION:-}"

#################################################################################
# Functions
#################################################################################

function usage() {
    cat <<-EOF
Usage:
  generate-builder-config.sh [OPTIONS]

The following OPTIONS are accepted:
-${OPT_SHORT_HELP}, --${OPT_LONG_HELP}              Print this usage information message
-${OPT_SHORT_INPUT}, --${OPT_LONG_INPUT}             Path to the config template (default: ${DEFAULT_INPUT_PATH})
-${OPT_SHORT_OUTPUT}, --${OPT_LONG_OUTPUT}            Path to output the generated config to (default: ${DEFAULT_OUTPUT_PATH})

The following ENVIRONMENT VARIABLES are accepted:
DIST_VERSION            The version to embed into the collector binary (default: ${DEFAULT_DIST_VERSION})
OT_VERSION              The version of opentelemetry-collector-contrib to use (required)
EOF
}

function exit_with_usage() {
    usage >&2
    exit 1
}

function exit_with_usage_and_error() {
    local msg="$1"
    readonly msg

    usage >&2
    printf "\nError: %s\n" "${msg}" >&2
    exit 1
}

function parse_option_with_argument() {
    if [[ -n "$2" && "$2" != -* ]]; then
        echo "$2"
    else
        exit_with_usage_and_error "$1 requires a non-empty argument"
    fi
}

function parse_options() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
        "-${OPT_SHORT_INPUT}" | "--${OPT_LONG_INPUT}")
            INPUT_PATH="$(parse_option_with_argument "$@")"
            shift 2
            ;;
        "-${OPT_SHORT_OUTPUT}" | "--${OPT_LONG_OUTPUT}")
            OUTPUT_PATH="$(parse_option_with_argument "$@")"
            shift 2
            ;;
        "-${OPT_SHORT_HELP}" | "--${OPT_LONG_HELP}")
            exit_with_usage
            ;;
        -*)
            exit_with_usage_and_error "unrecognized option: $1"
            ;;
        *)
            # Positional argument
            exit_with_usage_and_error "position arguments are unsupported: $1"
            ;;
        esac
    done
}

function missing_environment_variable() {
    exit_with_usage_and_error "missing required environment variable: $1"
}

parse_options "$@"

# Verify that the input template exists
if [ ! -f "${INPUT_PATH}" ]; then
    exit_with_usage_and_error "file not found: ${INPUT_PATH}"
fi

# Verify that the required environment variables are set, and then export them
if [ "${DIST_VERSION}" == "" ]; then
    missing_environment_variable "DIST_VERSION"
fi
export DIST_VERSION

if [ "${OT_VERSION}" == "" ]; then
    missing_environment_variable "OT_VERSION"
fi
export OT_VERSION

# Replace environment variables in the template and output config file
yq -e '(.. | select(tag == "!!str")) |= envsubst(nu)' "${INPUT_PATH}" >"${OUTPUT_PATH}"

cat "${OUTPUT_PATH}"
