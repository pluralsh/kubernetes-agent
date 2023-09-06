"""
Macros for cmd.
"""

load("@io_bazel_rules_go//go:def.bzl", "go_binary")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_image_index", "oci_push")
load("@rules_pkg//:pkg.bzl", "pkg_tar")

_x_defs = {
    "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd.Version": "{STABLE_BUILD_GIT_TAG}",
    "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd.Commit": "{STABLE_BUILD_GIT_COMMIT}",
    "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd.BuildTime": "{BUILD_TIME}",
}

# This can be overridden using command line flags. See https://docs.aspect.build/rules/rules_oci/docs/push.
_agentk_repo = "registry.gitlab.com/gitlab-org/cluster-integration/gitlab-agent"

def define_command_targets(name, binary_embed, arm_targets = True, arm64_targets = True):
    go_binary(
        name = name,
        embed = binary_embed,
        visibility = ["//visibility:public"],
        x_defs = _x_defs,
    )

    go_binary(
        name = "%s_race" % name,
        embed = binary_embed,
        race = "on",
        tags = ["manual"],
        visibility = ["//visibility:public"],
        x_defs = _x_defs,
    )

    images = []
    debug_images = []
    images.append(binary_and_image(
        name = name,
        arch = "amd64",
        arch_variant = None,
        debug = False,
        race = "off",
        binary_embed = binary_embed,
    ))

    debug_images.append(binary_and_image(
        name = name,
        arch = "amd64",
        arch_variant = None,
        debug = True,
        race = "on",
        binary_embed = binary_embed,
    ))

    if arm_targets:
        images.append(binary_and_image(
            name = name,
            arch = "arm",
            arch_variant = "v7",
            debug = False,
            race = "off",
            binary_embed = binary_embed,
        ))
        debug_images.append(binary_and_image(
            name = name,
            arch = "arm",
            arch_variant = "v7",
            debug = True,
            race = "off",
            binary_embed = binary_embed,
        ))

    if arm64_targets:
        images.append(binary_and_image(
            name = name,
            arch = "arm64",
            arch_variant = "v8",
            debug = False,
            race = "off",
            binary_embed = binary_embed,
        ))
        debug_images.append(binary_and_image(
            name = name,
            arch = "arm64",
            arch_variant = "v8",
            debug = True,
            race = "off",
            binary_embed = binary_embed,
        ))

    oci_image_index(
        name = "%s_index" % name,
        images = images,
        visibility = ["//visibility:public"],
        tags = ["manual"],
    )

    oci_image_index(
        name = "%s_index_debug" % name,
        images = debug_images,
        visibility = ["//visibility:public"],
        tags = ["manual"],
    )

    oci_push(
        name = "push",
        image = ":%s_index" % name,
        repository = _agentk_repo,
        tags = ["manual"],
    )

    oci_push(
        name = "push_debug",
        image = ":%s_index_debug" % name,
        repository = _agentk_repo,
        tags = ["manual"],
    )

def binary_and_image(name, arch, arch_variant, debug, race, binary_embed):
    binary_arch = arch if arch_variant == None else "%s_%s" % (arch, arch_variant)
    if debug:
        binary_name = "%s_linux_%s_debug" % (name, binary_arch)
        base = "@distroless_base_debug_nonroot_linux_%s" % binary_arch
    else:
        binary_name = "%s_linux_%s" % (name, binary_arch)
        base = "@distroless_static_nonroot_linux_%s" % binary_arch
    go_binary(
        name = binary_name,
        embed = binary_embed,
        goarch = arch,
        goos = "linux",
        race = race,
        tags = ["manual"],
        visibility = ["//visibility:public"],
        x_defs = _x_defs,
    )

    tar_name = "%s_tar" % binary_name
    pkg_tar(
        name = tar_name,
        srcs = [":" + binary_name],
        tags = ["manual"],
    )

    oci_image_name = "%s_image" % binary_name
    oci_image(
        name = oci_image_name,
        base = base,
        tars = [":" + tar_name],
        tags = ["manual"],
        visibility = ["//visibility:public"],
        entrypoint = ["/" + binary_name],
    )

    return ":" + oci_image_name
