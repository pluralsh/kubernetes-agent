"""
Macros for cmd.
"""

load("@io_bazel_rules_go//go:def.bzl", "go_binary")
load("@io_bazel_rules_docker//container:container.bzl", "container_bundle")
load("@io_bazel_rules_docker//contrib:push-all.bzl", "container_push")
load("@io_bazel_rules_docker//lang:image.bzl", "app_layer")

def push_bundle(name, images):
    bundle_name = name + "_bundle"
    container_bundle(
        name = bundle_name,
        images = images,
        tags = ["manual"],
        visibility = ["//visibility:public"],
    )
    container_push(
        name = name,
        bundle = ":" + bundle_name,
        format = "Docker",
        tags = ["manual"],
        visibility = ["//visibility:public"],
    )

def define_command_targets(
        name,
        binary_embed,
        race_targets = True,
        arm_targets = True,
        arm64_targets = True,
        base_image = "@go_image_static//image",
        base_image_arm = "@go_image_static_arm//image",
        base_image_arm64 = "@go_image_static_arm64//image",
        base_image_race = "@go_debug_image_base//image",
        base_image_arm_race = "@go_debug_image_base_arm//image",
        base_image_arm64_race = "@go_debug_image_base_arm64//image"):
    x_defs = {
        "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd.Version": "{STABLE_BUILD_GIT_TAG}",
        "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd.Commit": "{STABLE_BUILD_GIT_COMMIT}",
        "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd.BuildTime": "{BUILD_TIME}",
    }
    go_binary(
        name = name,
        embed = binary_embed,
        visibility = ["//visibility:public"],
        x_defs = x_defs,
    )

    go_binary(
        name = "%s_linux" % name,
        embed = binary_embed,
        goarch = "amd64",
        goos = "linux",
        tags = ["manual"],
        visibility = ["//visibility:public"],
        x_defs = x_defs,
    )

    go_image(
        name = "container",
        base = base_image,
        binary = ":%s_linux" % name,
        tags = ["manual"],
        visibility = ["//visibility:public"],
    )

    if arm_targets:
        go_binary(
            name = "%s_linux_arm" % name,
            embed = binary_embed,
            goarch = "arm",
            goos = "linux",
            tags = ["manual"],
            visibility = ["//visibility:public"],
            x_defs = x_defs,
        )

        go_image(
            name = "container_arm",
            base = base_image_arm,
            binary = ":%s_linux_arm" % name,
            architecture = "arm",
            tags = ["manual"],
            visibility = ["//visibility:public"],
        )

    if arm64_targets:
        go_binary(
            name = "%s_linux_arm64" % name,
            embed = binary_embed,
            goarch = "arm64",
            goos = "linux",
            tags = ["manual"],
            visibility = ["//visibility:public"],
            x_defs = x_defs,
        )

        go_image(
            name = "container_arm64",
            base = base_image_arm64,
            binary = ":%s_linux_arm64" % name,
            architecture = "arm64",
            tags = ["manual"],
            visibility = ["//visibility:public"],
        )

    if race_targets:
        go_binary(
            name = "%s_race" % name,
            embed = binary_embed,
            race = "on",
            tags = ["manual"],
            visibility = ["//visibility:public"],
            x_defs = x_defs,
        )

        go_binary(
            name = "%s_linux_race" % name,
            embed = binary_embed,
            goarch = "amd64",
            goos = "linux",
            race = "on",
            tags = ["manual"],
            visibility = ["//visibility:public"],
            x_defs = x_defs,
        )

        go_image(
            name = "container_race",
            base = base_image_race,
            binary = ":%s_linux_race" % name,
            tags = ["manual"],
            visibility = ["//visibility:public"],
        )

    if race_targets and arm_targets:
        # Will only work on an arm machine because otherwise cross compilation with CGO requires a properly setup crosstool.
        #        go_binary(
        #            name = "%s_linux_arm_race" % name,
        #            embed = binary_embed,
        #            goarch = "arm",
        #            goos = "linux",
        #            race = "on",
        #            tags = ["manual"],
        #            visibility = ["//visibility:public"],
        #            x_defs = x_defs,
        #        )
        #
        go_image(
            name = "container_arm_race",
            base = base_image_arm_race,
            binary = ":%s_linux_arm" % name,  # not a race binary, but this image is still good for debugging
            architecture = "arm",
            tags = ["manual"],
            visibility = ["//visibility:public"],
        )

    if race_targets and arm64_targets:
        # Will only work on an arm64 machine because otherwise cross compilation with CGO requires a properly setup crosstool.
        #        go_binary(
        #            name = "%s_linux_arm64_race" % name,
        #            embed = binary_embed,
        #            goarch = "arm64",
        #            goos = "linux",
        #            race = "on",
        #            tags = ["manual"],
        #            visibility = ["//visibility:public"],
        #            x_defs = x_defs,
        #        )
        #
        go_image(
            name = "container_arm64_race",
            base = base_image_arm64_race,
            binary = ":%s_linux_arm64" % name,  # not a race binary, but this image is still good for debugging
            architecture = "arm64",
            tags = ["manual"],
            visibility = ["//visibility:public"],
        )

# This is a fork of load("@io_bazel_rules_docker//go:image.bzl", "go_image")
# to pick up "architecture" support from https://github.com/bazelbuild/rules_docker/pull/1936.
def go_image(name, base, deps = [], layers = [], binary = None, architecture = None, **kwargs):
    """Constructs a container image wrapping a go_binary target.

  Args:
    name: Name of the go_image target.
    base: Base image to use to build the go_image.
    deps: Dependencies of the go image target.
    binary: An alternative binary target to use instead of generating one.
    layers: Augments "deps" with dependencies that should be put into their own layers.
    **kwargs: See go_binary.
  """
    if layers:
        print("go_image does not benefit from layers=[], got: %s" % layers)

    if not binary:
        binary = name + ".binary"
        go_binary(name = binary, deps = deps + layers, **kwargs)
    elif deps:
        fail("kwarg does nothing when binary is specified", "deps")

    #    if not base:
    #        base = STATIC_DEFAULT_BASE if kwargs.get("pure") == "on" else DEFAULT_BASE

    tags = kwargs.get("tags", None)
    for index, dep in enumerate(layers):
        base = app_layer(name = "%s.%d" % (name, index), base = base, dep = dep, tags = tags)
        base = app_layer(name = "%s.%d-symlinks" % (name, index), base = base, dep = dep, binary = binary, tags = tags)

    visibility = kwargs.get("visibility", None)
    restricted_to = kwargs.get("restricted_to", None)
    compatible_with = kwargs.get("compatible_with", None)
    app_layer(
        name = name,
        architecture = architecture,  # added in this "fork"
        base = base,
        binary = binary,
        visibility = visibility,
        tags = tags,
        args = kwargs.get("args"),
        data = kwargs.get("data"),
        testonly = kwargs.get("testonly"),
        restricted_to = restricted_to,
        compatible_with = compatible_with,
    )
