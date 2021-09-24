load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_proto_grpc//go:defs.bzl", "go_grpc_compile", "go_proto_compile", "go_validate_compile")
load("//build:build.bzl", "copy_to_workspace")
load("@bazel_skylib//lib:paths.bzl", "paths")
load("@com_github_ash2k_bazel_tools//multirun:def.bzl", "multirun")

def go_proto_generate(src, workspace_relative_target_directory, deps = []):
    extract_targets = _go_proto_copy(src, ":proto_compile", workspace_relative_target_directory, deps)
    go_proto_compile(
        name = "proto_compile",
        tags = ["manual"],
        protos = [":proto"],
    )
    create_extract_target(extract_targets)

def create_extract_target(extract_targets):
    if len(extract_targets) == 1:
        native.alias(
            name = "extract_generated",
            actual = extract_targets[0],
            visibility = ["//visibility:public"],
        )
    else:
        multirun(
            name = "extract_generated",
            commands = extract_targets,
            visibility = ["//visibility:public"],
        )

def _go_proto_copy(src, go_compile_label, workspace_relative_target_directory, deps):
    proto_library(
        name = "proto",
        srcs = [src],
        tags = ["manual"],
        visibility = ["//visibility:public"],
        deps = deps,
    )

    copy_to_workspace(
        name = "extract_generated_proto",
        file_to_copy = paths.split_extension(src)[0] + ".pb.go",
        label = go_compile_label,
        workspace_relative_target_directory = workspace_relative_target_directory,
    )

    extract_targets = [":extract_generated_proto"]

    if "@com_github_envoyproxy_protoc_gen_validate//validate:validate_proto" in deps:
        # https://github.com/rules-proto-grpc/rules_proto_grpc/issues/146
        go_validate_compile(
            name = "proto_compile_validate",
            tags = ["manual"],
            protos = [":proto"],
        )
        copy_to_workspace(
            name = "extract_generated_validator",
            file_to_copy = paths.split_extension(src)[0] + ".pb.validate.go",
            label = ":proto_compile_validate",
            workspace_relative_target_directory = workspace_relative_target_directory,
        )
        extract_targets.append(":extract_generated_validator")

    return extract_targets

def go_grpc_generate(src, workspace_relative_target_directory, deps = []):
    extract_targets = _go_proto_copy(src, ":grpc_compile", workspace_relative_target_directory, deps)
    go_grpc_compile(
        name = "grpc_compile",
        tags = ["manual"],
        protos = [":proto"],
    )

    copy_to_workspace(
        name = "extract_generated_grpc",
        file_to_copy = paths.split_extension(src)[0] + "_grpc.pb.go",
        label = ":grpc_compile",
        workspace_relative_target_directory = workspace_relative_target_directory,
    )

    extract_targets.append(":extract_generated_grpc")
    create_extract_target(extract_targets)
