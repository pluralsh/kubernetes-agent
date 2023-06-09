load("@io_bazel_rules_go//go:def.bzl", "go_test")

# go_custom_test is a macro around go_test that sets size="small" and race="on" if these
# arguments are not set explicitly.
def go_custom_test(size = "small", race = "on", **kwargs):
    go_test(
        size = size,
        race = race,
        **kwargs
    )
