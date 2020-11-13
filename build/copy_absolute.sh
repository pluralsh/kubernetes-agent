#!/usr/bin/env bash

set -e -o pipefail

# This is a helper script that allows copying files, generated by a bazel rule, into a directory, specified
# by an absolute path.

if (( $# != 3 )); then
  echo 'Not enough or too many command line arguments' >&2
  exit 1
fi

source_files="$1"
file_to_copy="$2"
absolute_target_directory="$3"

# Don't want to double quote because $(rootpaths //label) in build.bzl expands into a single argument which
# contains space-separated file names.
# shellcheck disable=SC2068
for file in $source_files
do
  name=$(basename "$file")
  if [[ "$name" == "$file_to_copy" ]]
  then
    to="$absolute_target_directory/$name"
    cp "$file" "$to"
    chmod +w "$to"
    break
  fi
done
