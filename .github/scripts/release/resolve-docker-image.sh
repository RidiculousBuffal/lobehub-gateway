#!/usr/bin/env bash
set -euo pipefail

# Resolve the target image and login registry from the optional DOCKER_IMAGE
# secret. A first path segment without a dot or colon is a Docker Hub
# namespace, not a registry host.
if [[ -n "${DOCKER_IMAGE:-}" ]]; then
  IMAGE="$DOCKER_IMAGE"
  if [[ "$IMAGE" == */* ]]; then
    REGISTRY="${IMAGE%%/*}"
    if [[ ! "$REGISTRY" =~ \. ]] && [[ ! "$REGISTRY" =~ : ]]; then
      REGISTRY="docker.io"
    fi
  else
    REGISTRY="docker.io"
  fi
else
  IMAGE="ghcr.io/$(echo "$GITHUB_REPOSITORY" | tr '[:upper:]' '[:lower:]')"
  REGISTRY="ghcr.io"
fi

{
  echo "image=$IMAGE"
  echo "registry=$REGISTRY"
} >> "$GITHUB_OUTPUT"
