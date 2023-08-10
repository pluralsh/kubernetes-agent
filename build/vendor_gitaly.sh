#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

A="${AGENT_DIR:-"$HOME/src/gitlab-agent/internal/gitaly/vendored"}"
G="${GITALY_DIR:-"$HOME/src/gitaly"}"

rm -rf "$A/backoff"
cp -R "$G/internal/backoff" "$A"
rm "$A/backoff/"*_test.go

rm -rf "$A/structerr"
mkdir -p "$A/structerr"
cp -R "$G/internal/structerr/error.go" "$A/structerr"

rm -rf "$A/dnsresolver"
cp -R "$G/internal/grpc/dnsresolver" "$A/dnsresolver"
rm "$A/dnsresolver/"*_test.go

rm -rf "$A/grpc/client"
cp -R "$G/internal/grpc/client" "$A/grpc/client"
rm "$A/grpc/client/"*_test.go

rm -rf "$A/client"
cp -R "$G/client" "$A"
rm -rf "$A/client/"*_test.go "$A/client/testdata" "$A/client/receive_pack.go" "$A/client/sidechannel.go" "$A/client/upload_archive.go" "$A/client/upload_pack.go"

rm -rf "$A/gitalyauth"
cp -R "$G/auth" "$A/gitalyauth"
rm "$A/gitalyauth/"*_test.go "$A/gitalyauth/README.md"

mkdir -p "$A/gitalypb"
rm "$A/gitalypb/"*.proto || true

# Only copy what we need
for FILE in lint errors shared commit service_config smarthttp packfile
do
  cp "$G/proto/$FILE.proto" "$A/gitalypb"
done

GP='gitlab.com/gitlab-org/gitaly/v16'
AV='gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored'

sed -i '' -e "s|\"$GP/internal/backoff\"|\"$AV/backoff\"|g" -- "$A/dnsresolver/"*.go "$A/client/"*.go
sed -i '' -e "s|\"$GP/internal/structerr\"|\"$AV/structerr\"|g" -- "$A/dnsresolver/"*.go
sed -i '' -e "s|\"$GP/internal/grpc/dnsresolver\"|\"$AV/dnsresolver\"|g" -- "$A/client/"*.go "$A/grpc/client/"*.go
sed -i '' -e "s|\"$GP/internal/grpc/client\"|\"$AV/grpc/client\"|g" -- "$A/client/"*.go
sed -i '' -e "s|\"$GP/auth\"|\"$AV/gitalyauth\"|g" -- "$A/client/"*.go

sed -i '' -e "s|gitalyx509 \"$GP/internal/x509\"|\"crypto/x509\"|g" -- "$A/grpc/client/"*.go
sed -i '' -e "s|gitalyx509\.SystemCertPool()|x509.SystemCertPool()|g" -- "$A/grpc/client/"*.go

sed -i '' -e "s|\"$GP/proto/go/gitalypb\"|\"$AV/gitalypb\"|g" -- "$A/gitalypb/"*.proto "$A/grpc/client/"*.go
