#!/bin/sh
#
# Run `golangci-lint` as per the `treefmt` specification.
#
# Basically, if any go file has changed, run `golangci-lint` on the
# full project.

set -eu

golangci-lint run --fix ./...
