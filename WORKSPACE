workspace(name = "gitlab_k8s_agent")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# When updating rules_go make sure to update org_golang_x_tools dependency below by copying it from
# https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
# Also update to the same version/commit in go.mod.
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "099a9fb96a376ccbbb7d291ed4ecbdfd42f6bc822ab77ae6f1b5cb9e914e94fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "efbbba6ac1a4fd342d5122cbdfdb82aeb2cf2862e35022c752eaddffada7c3f3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
    ],
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "b1e80761a8a8243d03ebca8845e9cc1ba6c82ce7c5179ce2b295cd36f7e394bf",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.25.0/rules_docker-v0.25.0.tar.gz"],
)

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "e3bb0dc8b0274ea1aca75f1f8c0c835adbe589708ea89bf698069d0790701ea3",
    strip_prefix = "buildtools-5.1.0",
    urls = ["https://github.com/bazelbuild/buildtools/archive/5.1.0.tar.gz"],
)

http_archive(
    name = "com_github_ash2k_bazel_tools",
    sha256 = "f12cdb947d8c92c7bbed24f4f4492a23b9b1cf7f384d3662d99ee3753d14c15a",
    strip_prefix = "bazel-tools-4daedde3ec61a03db841c8a9ca68288972e25a82",
    urls = ["https://github.com/ash2k/bazel-tools/archive/4daedde3ec61a03db841c8a9ca68288972e25a82.tar.gz"],
)

http_archive(
    name = "rules_proto",
    sha256 = "80d3a4ec17354cccc898bfe32118edd934f851b03029d63ef3fc7c8663a7415c",
    strip_prefix = "rules_proto-5.3.0-21.5",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/5.3.0-21.5.tar.gz",
    ],
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "bbe4db93499f5c9414926e46f9e35016999a4e9f6e3522482d3760dc61011070",
    strip_prefix = "rules_proto_grpc-4.2.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/4.2.0.tar.gz"],
)

http_archive(
    name = "bazel_skylib",
    sha256 = "f7be3474d42aae265405a592bb7da8e171919d74c16f082a5457840f06054728",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.2.1/bazel-skylib-1.2.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.2.1/bazel-skylib-1.2.1.tar.gz",
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
    sum = "h1:Y7DTJMR6zs1xkS/upamJYk0SxxN4C9AqRd77jmZnyY4=",
    version = "v1.10.0",
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
    sum = "h1:pDDYmo0QadUPal5fwXoY1pmMpFcdyhXOmL5drCrI3vU=",
    version = "v1.10.0",
)

# See https://github.com/open-telemetry/opentelemetry-go-contrib/issues/872
go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracehttp",
    build_directives = [
        "gazelle:resolve go go.opentelemetry.io/otel/exporters/otlp/internal @io_opentelemetry_go_otel//exporters/otlp/internal",
    ],
    build_file_proto_mode = "disable_global",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp",
    sum = "h1:S8DedULB3gp93Rh+9Z+7NTEv+6Id/KYS7LDyipZ9iCE=",
    version = "v1.10.0",
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
    sum = "h1:B2cR/FAaiMtYDHv5BQpaqtkjGuWQIgr2KQZtHQ7f6i8=",
    version = "v0.6.8",
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
    sha256 = "4e3d94e7bf8dde5dad681c5ddddda6e634f8c8c500683fdf3d2f77a9c086702d",
    strip_prefix = "tools-0.1.12",
    # v0.1.12, latest as of 2022-09-10
    urls = [
        "https://mirror.bazel.build/github.com/golang/tools/archive/refs/tags/v0.1.12.zip",
        "https://github.com/golang/tools/archive/refs/tags/v0.1.12.zip",
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
    version = "1.18.6",
)

gazelle_dependencies()

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

# Latest images as of 2022-06-18

# debug-nonroot-amd64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian11
container_pull(
    name = "go_debug_image_base",
    digest = "sha256:2d418d1fdbc67131ed95941eab4b8fc6192487a33b54ed91cdb543c375f30486",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
)

# nonroot-amd64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static-debian11
container_pull(
    name = "go_image_static",
    digest = "sha256:c97d534e9f2b9bd834813b0e676fb2ac53c55719a83570c65b64f8319b4de6c9",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
)

# debug-nonroot-arm from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian11
container_pull(
    name = "go_debug_image_base_arm",
    architecture = "arm",
    digest = "sha256:8f12d5312f92b871fb0d778c7cc3aa364d2de910eb1d17ddc62509b4cacfbe5a",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
)

# nonroot-arm from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static-debian11
container_pull(
    name = "go_image_static_arm",
    architecture = "arm",
    digest = "sha256:ec29fb858b71ccd29256de0b5be8022f5300eba0d2499205d12b09215634ddd1",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
)

# debug-nonroot-arm64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/base-debian11
container_pull(
    name = "go_debug_image_base_arm64",
    architecture = "arm64",
    digest = "sha256:7eb3d3112f7c1dd704d530dc85b01b169633260a0b06310b536b3fb3499c46d3",
    registry = "gcr.io",
    repository = "distroless/base-debian11",
)

# nonroot-arm64 from https://console.cloud.google.com/gcr/images/distroless/GLOBAL/static-debian11
container_pull(
    name = "go_image_static_arm64",
    architecture = "arm64",
    digest = "sha256:00d273f8a59b3e0fa9d244304cc960952e00ab5f3e43d37143f72091d7641c79",
    registry = "gcr.io",
    repository = "distroless/static-debian11",
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")
load("@com_github_ash2k_bazel_tools//buildozer:deps.bzl", "buildozer_dependencies")
load("@com_github_ash2k_bazel_tools//multirun:deps.bzl", "multirun_dependencies")
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")
load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")
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
