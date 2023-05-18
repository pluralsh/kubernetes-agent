load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_proto_grpc//go:defs.bzl", "go_proto_compile")
load("@rules_proto_grpc//doc:defs.bzl", "doc_markdown_compile")
load("@bazel_skylib//lib:paths.bzl", "paths")
load("@aspect_bazel_lib//lib:write_source_files.bzl", "write_source_files")
load("//build:proto_def.bzl", "go_grpc_compile", "go_validate_compile")

def _common(src, deps):
    proto_library(
        name = "proto",
        srcs = [src],
        tags = ["manual"],
        visibility = ["//visibility:public"],
        deps = deps,
    )
    name_only = paths.split_extension(src)[0]
    extract_targets = {}
    go_proto_compile(
        name = "proto_compile",
        tags = ["manual"],
        protos = [":proto"],
    )
    extract_targets[name_only + ".pb.go"] = ":proto_compile"

    if "@com_github_envoyproxy_protoc_gen_validate//validate:validate_proto" in deps:
        go_validate_compile(
            name = "proto_compile_validate",
            tags = ["manual"],
            protos = [":proto"],
        )
        extract_targets[name_only + ".pb.validate.go"] = ":proto_compile_validate"

    return extract_targets

def _create_extract_target(extract_targets):
    write_source_files(
        name = "extract_generated",
        files = extract_targets,
        diff_test = False,
        tags = ["manual"],
        visibility = ["//visibility:public"],
    )

def go_proto_generate(src, deps = [], with_md_docs = False):
    extract_targets = _common(src, deps)

    if with_md_docs:
        filename = paths.split_extension(src)[0] + "_proto_docs"
        doc_markdown_compile(
            name = filename,
            protos = [":proto"],
        )
        extract_targets[filename + ".md"] = ":" + filename

    _create_extract_target(extract_targets)

def go_grpc_generate(src, deps = []):
    extract_targets = _common(src, deps)

    go_grpc_compile(
        name = "grpc_compile",
        tags = ["manual"],
        protos = [":proto"],
    )
    extract_targets[paths.split_extension(src)[0] + "_grpc.pb.go"] = ":grpc_compile"

    _create_extract_target(extract_targets)
