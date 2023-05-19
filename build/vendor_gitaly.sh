#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

A="$HOME/src/gitlab-agent/internal/gitaly/vendored"
G="$HOME/src/gitaly"

rm -rf "$A/backoff"
cp -R "$G/internal/backoff" "$A"
rm "$A/backoff/"*_test.go

rm -rf "$A/structerr"
mkdir -p "$A/structerr"
cp -R "$G/internal/structerr/error.go" "$A/structerr"

rm -rf "$A/dnsresolver"
cp -R "$G/internal/grpc/dnsresolver" "$A/dnsresolver"
rm "$A/dnsresolver/"*_test.go

rm -rf "$A/internal_client"
cp -R "$G/internal/gitaly/client" "$A/internal_client"
rm "$A/internal_client/"*_test.go

rm -rf "$A/client"
cp -R "$G/client" "$A"
rm -rf "$A/client/"*_test.go "$A/client/testdata" "$A/client/receive_pack.go" "$A/client/sidechannel.go" "$A/client/upload_archive.go" "$A/client/upload_pack.go"

rm -rf "$A/gitalyauth"
cp -R "$G/auth" "$A/gitalyauth"
rm "$A/gitalyauth/"*_test.go "$A/gitalyauth/README.md"

mkdir -p "$A/gitalypb"
rm "$A/gitalypb/"*.proto || true

# Only copy what we need
for FILE in lint errors shared commit service_config smarthttp
do
  cp "$G/proto/$FILE.proto" "$A/gitalypb"
done

GP='gitlab.com/gitlab-org/gitaly/v16'
AV='gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored'

sed -i '' -e "s|\"$GP/internal/backoff\"|\"$AV/backoff\"|g" -- "$A/dnsresolver/"*.go "$A/client/"*.go
sed -i '' -e "s|\"$GP/internal/structerr\"|\"$AV/structerr\"|g" -- "$A/dnsresolver/"*.go
sed -i '' -e "s|\"$GP/internal/grpc/dnsresolver\"|\"$AV/dnsresolver\"|g" -- "$A/client/"*.go "$A/internal_client/"*.go
sed -i '' -e "s|\"$GP/internal/gitaly/client\"|\"$AV/internal_client\"|g" -- "$A/client/"*.go
sed -i '' -e "s|\"$GP/auth\"|\"$AV/gitalyauth\"|g" -- "$A/client/"*.go

sed -i '' -e "s|gitalyx509 \"$GP/internal/x509\"|\"crypto/x509\"|g" -- "$A/internal_client/"*.go
sed -i '' -e "s|gitalyx509\.SystemCertPool()|x509.SystemCertPool()|g" -- "$A/internal_client/"*.go

sed -i '' -e "s|\"$GP/proto/go/gitalypb\"|\"$AV/gitalypb\"|g" -- "$A/gitalypb/"*.proto "$A/internal_client/"*.go
