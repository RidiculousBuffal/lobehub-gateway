#!/usr/bin/env bash
set -euo pipefail

if [[ "$INPUT_VERSION" == v* ]]; then
  echo "::error::Release tags do not use a v prefix. Use 0.1.0 instead of v0.1.0"
  exit 1
fi

SEMVER_REGEX='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(alpha|beta|rc)\.(0|[1-9][0-9]*))?$'
if [[ ! "$INPUT_VERSION" =~ $SEMVER_REGEX ]]; then
  echo "::error::Version must be in the form 0.1.0 or 0.1.0-rc.1. Pre-release tags allow only alpha, beta, rc, with no build metadata"
  exit 1
fi

case "$INPUT_LATEST" in
  true | false) ;;
  *)
    echo "::error::latest must be true or false, got $INPUT_LATEST"
    exit 1
    ;;
esac

TARGET_SHA="$(git rev-parse --verify HEAD)"
if [[ "$INPUT_VERSION" == *-* ]]; then
  PRERELEASE=true
else
  PRERELEASE=false
fi

if [[ "$PRERELEASE" == "true" && "$INPUT_LATEST" == "true" ]]; then
  echo "::error::Pre-release versions cannot be published as latest. Disable latest, or use a stable version number"
  exit 1
fi

if [[ "$INPUT_LATEST" == "true" ]]; then
  MASTER_SHA="$(git ls-remote origin refs/heads/master | cut -f1)"
  if [[ "$TARGET_SHA" != "$MASTER_SHA" ]]; then
    echo "::error::latest can only be published from the current HEAD of origin/master. workflow_ref=$GITHUB_REF_NAME target_sha=$TARGET_SHA origin/master=$MASTER_SHA"
    exit 1
  fi
fi

{
  echo "version=$INPUT_VERSION"
  echo "target_sha=$TARGET_SHA"
  echo "prerelease=$PRERELEASE"
  echo "publish_latest=$INPUT_LATEST"
} >> "$GITHUB_OUTPUT"
