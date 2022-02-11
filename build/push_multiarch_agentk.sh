#!/usr/bin/env bash

set -e
basetag="$1"

if [[ -z "$CI_REGISTRY_IMAGE" ]]; then
    echo '$CI_REGISTRY_IMAGE must be present'
    return 1
fi

if [[ -z "$basetag" ]]; then
    echo '$0 takes exactly one argument'
    return 1
fi

docker manifest create "${CI_REGISTRY_IMAGE}/agentk:${basetag}" \
    --amend "${CI_REGISTRY_IMAGE}/agentk:${basetag}-amd64" \
    --amend "${CI_REGISTRY_IMAGE}/agentk:${basetag}-arm64"

docker manifest push "${CI_REGISTRY_IMAGE}/agentk:${basetag}"

docker manifest create "${CI_REGISTRY_IMAGE}/agentk:${basetag}-race" \
    --amend "${CI_REGISTRY_IMAGE}/agentk:${basetag}-amd64-race" \
    --amend "${CI_REGISTRY_IMAGE}/agentk:${basetag}-arm64-race"

docker manifest push "${CI_REGISTRY_IMAGE}/agentk:${basetag}-race"
