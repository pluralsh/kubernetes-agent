#!/usr/bin/env bash

# This command is used by bazel as the workspace_status_command
# to implement build stamping with git information.

set -o errexit
set -o nounset
set -o pipefail

# If the GIT_COMMIT or GIT_TAG variables are not already set
# then it will not try and use git commands to set them.
# See: https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/253
[ -z "${GIT_COMMIT:-}" ] && GIT_COMMIT=$(git rev-parse --short HEAD)
[ -z "${GIT_TAG:-}" ] && GIT_TAG=$(git tag --points-at HEAD 2>/dev/null || true)
GIT_TAG="${GIT_TAG:="v0.0.0"}"

BUILD_TIME=$(date -u +%Y%m%d.%H%M%S)
# Prefix with STABLE_ so that these values are saved to stable-status.txt
# instead of volatile-status.txt.
# Stamped rules will be retriggered by changes to stable-status.txt, but not by
# changes to volatile-status.txt.
# See https://docs.bazel.build/versions/master/user-manual.html#flag--workspace_status_command
cat <<EOF
STABLE_BUILD_GIT_COMMIT ${GIT_COMMIT-}
STABLE_BUILD_GIT_TAG ${GIT_TAG-}
BUILD_TIME ${BUILD_TIME-}
EOF
