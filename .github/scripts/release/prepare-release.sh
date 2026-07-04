#!/usr/bin/env bash
set -euo pipefail

SEMVER_REGEX='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(alpha|beta|rc)\.(0|[1-9][0-9]*))?$'
if [[ ! "$INPUT_VERSION" =~ $SEMVER_REGEX ]]; then
  echo "::error::Version must be in the form 0.1.0 or 0.1.0-rc.1. Pre-release tags allow only alpha, beta, rc, with no build metadata"
  exit 1
fi

if [[ "$INPUT_VERSION" == v* ]]; then
  echo "::error::Release tags do not use a v prefix. Use 0.1.0 instead of v0.1.0"
  exit 1
fi

case "$INPUT_LATEST" in
  true | false) ;;
  *)
    echo "::error::latest must be true or false, got $INPUT_LATEST"
    exit 1
    ;;
esac

git fetch --force --tags origin
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

EXISTING_TAG_SHA="$(git rev-parse --verify --quiet "refs/tags/$INPUT_VERSION^{}" || true)"
if [[ -n "$EXISTING_TAG_SHA" ]]; then
  if [[ "$EXISTING_TAG_SHA" != "$TARGET_SHA" ]]; then
    echo "::error::tag $INPUT_VERSION already exists but points to $EXISTING_TAG_SHA, not target $TARGET_SHA"
    exit 1
  fi

  if RELEASE_DRAFT="$(gh release view "$INPUT_VERSION" --json isDraft --jq .isDraft 2>/dev/null)"; then
    if [[ "$RELEASE_DRAFT" != "true" ]]; then
      echo "::error::Release $INPUT_VERSION is already published, cannot re-run the same version release"
      exit 1
    fi
    echo "::notice::Found existing draft release for this version, will reuse it"
  else
    echo "::notice::Found existing tag without a release, will create a draft release"
  fi
fi

if [[ "$INPUT_LATEST" == "true" ]]; then
  git fetch --force origin main
  MAIN_SHA="$(git rev-parse --verify origin/main)"
  if [[ "$TARGET_SHA" != "$MAIN_SHA" ]]; then
    echo "::error::latest can only be published from the current HEAD of origin/main. workflow_ref=$GITHUB_REF_NAME target_sha=$TARGET_SHA origin/main=$MAIN_SHA"
    exit 1
  fi
fi

{
  echo "version=$INPUT_VERSION"
  echo "target_sha=$TARGET_SHA"
  echo "prerelease=$PRERELEASE"
  echo "publish_latest=$INPUT_LATEST"
} >> "$GITHUB_OUTPUT"
