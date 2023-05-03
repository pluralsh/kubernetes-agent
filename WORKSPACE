workspace(name = "gitlab_k8s_agent")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# When updating rules_go make sure to update org_golang_x_tools dependency below by copying it from
# https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
# Also update to the same version/commit in go.mod.
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6dc2da7ab4cf5d7bfc7c949776b1b7c733f05e56edc4bcd9022bb249d2e2a996",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.39.1/rules_go-v0.39.1.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "727f3e4edd96ea20c29e8c2ca9e8d2af724d8c7778e7923a854b2c80952bc405",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.30.0/bazel-gazelle-v0.30.0.tar.gz",
    ],
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "b1e80761a8a8243d03ebca8845e9cc1ba6c82ce7c5179ce2b295cd36f7e394bf",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.25.0/rules_docker-v0.25.0.tar.gz"],
)

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "977a0bd4593c8d4c8f45e056d181c35e48aa01ad4f8090bdb84f78dca42f47dc",
    strip_prefix = "buildtools-6.1.2",
    urls = ["https://github.com/bazelbuild/buildtools/archive/refs/tags/v6.1.2.tar.gz"],
)

http_archive(
    name = "com_github_ash2k_bazel_tools",
    sha256 = "a911dab6711bc12a00f02cc94b66ced7dc57650e382ebd4f17c9cdb8ec2cbd56",
    strip_prefix = "bazel-tools-2add5bb84c2837a82a44b57e83c7414247aed43a",
    urls = ["https://github.com/ash2k/bazel-tools/archive/2add5bb84c2837a82a44b57e83c7414247aed43a.tar.gz"],
)

http_archive(
    name = "rules_proto",
    sha256 = "dc3fb206a2cb3441b485eb1e423165b231235a1ea9b031b4433cf7bc1fa460dd",
    strip_prefix = "rules_proto-5.3.0-21.7",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/5.3.0-21.7.tar.gz",
    ],
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "fb7fc7a3c19a92b2f15ed7c4ffb2983e956625c1436f57a3430b897ba9864059",
    strip_prefix = "rules_proto_grpc-4.3.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/4.3.0.tar.gz"],
)

http_archive(
    name = "bazel_skylib",
    sha256 = "b8a1527901774180afc798aeb28c4634bdccf19c4d98e7bdd1ce79d1fe9aaad7",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.4.1/bazel-skylib-1.4.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.4.1/bazel-skylib-1.4.1.tar.gz",
    ],
)

http_archive(
    name = "bazelruby_rules_ruby",
    sha256 = "5035393cb5043d49ca9de78acb9e8c8622a193f6463a57ad02383a622b6dc663",
    strip_prefix = "rules_ruby-0.6.0",
    urls = [
        "https://github.com/bazelruby/rules_ruby/archive/v0.6.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

# See https://github.com/open-telemetry/opentelemetry-go-contrib/issues/872
go_repository(
    name = "io_opentelemetry_go_otel",
    build_directives = [
        "gazelle:go_visibility @io_opentelemetry_go_otel_metric//:__subpackages__",
    ],  # keep
    build_file_proto_mode = "disable",
    importpath = "go.opentelemetry.io/otel",
    sum = "h1:/79Huy8wbf5DnIPhemGB+zEPVwnN6fuQybr/SRXa6hM=",
    version = "v1.14.0",
)

# See https://github.com/open-telemetry/opentelemetry-go-contrib/issues/872
go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace",
    build_directives = [
        "gazelle:resolve go go.opentelemetry.io/otel/exporters/otlp/internal @io_opentelemetry_go_otel//exporters/otlp/internal",
        "gazelle:resolve go go.opentelemetry.io/otel/exporters/otlp/internal/envconfig @io_opentelemetry_go_otel//exporters/otlp/internal/envconfig",
    ],
    build_file_proto_mode = "disable",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace",
    sum = "h1:TKf2uAs2ueguzLaxOCBXNpHxfO/aC7PAdDsSH0IbeRQ=",
    version = "v1.14.0",
)

# See https://github.com/open-telemetry/opentelemetry-go-contrib/issues/872
go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracehttp",
    build_directives = [
        "gazelle:resolve go go.opentelemetry.io/otel/exporters/otlp/internal @io_opentelemetry_go_otel//exporters/otlp/internal",
    ],
    build_file_proto_mode = "disable_global",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
    sum = "h1:3jAYbRHQAqzLjd9I4tzxwJ8Pk/N6AqBcF6m1ZHrxG94=",
    version = "v1.14.0",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    build_file_proto_mode = "disable_global",
    build_naming_convention = "go_default_library",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    patch_args = ["-p1"],
    # patch addresses https://github.com/bazelbuild/bazel-gazelle/issues/941
    # patch created by manually editing the build file and running `diff -urN protoc-gen-validate protoc-gen-validate-copy`
    # or `diff -urN protoc-gen-validate/validate/BUILD protoc-gen-validate-copy/validate/BUILD` for a single file.
    patches = [
        "@gitlab_k8s_agent//build:validate_dependency.patch",
    ],
    sum = "h1:PS7VIOgmSVhWUEeZwTe7z7zouA22Cr590PzXKbZHOVY=",
    version = "v0.9.1",
)

# Copied from rules_go to keep patches in place
# https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
http_archive(
    name = "org_golang_x_tools",
    patch_args = ["-p1"],
    patches = [
        # deletegopls removes the gopls subdirectory. It contains a nested
        # module with additional dependencies. It's not needed by rules_go.
        # releaser:patch-cmd rm -rf gopls
        "@io_bazel_rules_go//third_party:org_golang_x_tools-deletegopls.patch",
        # releaser:patch-cmd gazelle -repo_root . -go_prefix golang.org/x/tools -go_naming_convention import_alias
        "@io_bazel_rules_go//third_party:org_golang_x_tools-gazelle.patch",
    ],
    sha256 = "9f20a20f29f4008d797a8be882ef82b69cf8f7f2b96dbdfe3814c57d8280fa4b",
    strip_prefix = "tools-0.7.0",
    # v0.7.0, latest as of 2023-03-27
    urls = [
        "https://github.com/golang/tools/archive/refs/tags/v0.7.0.zip",
        "https://mirror.bazel.build/github.com/golang/tools/archive/refs/tags/v0.7.0.zip",
    ],
)

# Here to set build_file_proto_mode=default. repositories.bzl sets it to disable_global which is not what we want.
go_repository(
    name = "com_github_lyft_protoc_gen_star",
    build_file_proto_mode = "default",
    importpath = "github.com/lyft/protoc-gen-star",
    sum = "h1:erE0rdztuaDq3bpGifD95wfoPrSZc95nGA6tbiNYh6M=",
    version = "v0.6.1",
)

load("//build:repositories.bzl", "go_repositories")

# gazelle:repository_macro build/repositories.bzl%go_repositories
go_repositories()

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

go_rules_dependencies()

go_register_toolchains(
    version = "1.19.8",
)

gazelle_dependencies()

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

# Images are managed by https://gitlab.com/gitlab-org/frontend/renovate-gitlab-bot/-/tree/main/renovate/projects/gitlab-agent.config.js
# DO NOT EDIT ================ START

# debug-nonroot-amd64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian11
container_pull(
    name = "go_debug_image_base",
    digest = "sha256:07feb0cbf22cc6873dd2ffed2c655afc499933696454399d8652dea60ef88cd6",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
)

# nonroot-amd64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static-debian11
container_pull(
    name = "go_image_static",
    digest = "sha256:81c9a17d330510c4c068d2570c2796cae06dc822014ddb79476ea136ca95ee71",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
)

# debug-nonroot-arm from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian11
container_pull(
    name = "go_debug_image_base_arm",
    architecture = "arm",
    digest = "sha256:3ddeb64148dea5b318e4d7585683c68fffe0364ba3b8e6c8accaf71dfb84259c",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
)

# nonroot-arm from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static-debian11
container_pull(
    name = "go_image_static_arm",
    architecture = "arm",
    digest = "sha256:8efffa2cf87ff83f16d8600e5cfbfbfc0e1776e6a56413ef7e0b515bfef23da2",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
)

# debug-nonroot-arm64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian11
container_pull(
    name = "go_debug_image_base_arm64",
    architecture = "arm64",
    digest = "sha256:987451edf11ded09813f6f101f31abbd68525f111f1a5f1e315904ffbf30b1b6",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
)

# nonroot-arm64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static-debian11
container_pull(
    name = "go_image_static_arm64",
    architecture = "arm64",
    digest = "sha256:42bf7118eb11d6e471f2e0740b8289452e5925c209da33447b00dda8f051a9ea",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
)

# DO NOT EDIT ================ END

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")
load("@com_github_ash2k_bazel_tools//buildozer:deps.bzl", "buildozer_dependencies")
load("@com_github_ash2k_bazel_tools//multirun:deps.bzl", "multirun_dependencies")
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")
load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")
load("@rules_proto_grpc//doc:repositories.bzl", rules_proto_grpc_doc_repos = "doc_repos")
load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)
load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

rules_proto_grpc_toolchains()

rules_proto_grpc_repos()

rules_proto_grpc_go_repos()

rules_proto_dependencies()

rules_proto_toolchains()

rules_proto_grpc_doc_repos()

container_repositories()

container_deps()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    go_image_repositories = "repositories",
)
load("@com_github_envoyproxy_protoc_gen_validate//:dependencies.bzl", pgv_third_party = "go_third_party")

go_image_repositories()

buildifier_dependencies()

buildozer_dependencies()

multirun_dependencies()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

pgv_third_party()
