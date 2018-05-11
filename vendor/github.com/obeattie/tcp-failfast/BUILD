load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "failfast.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "failfast_darwin.go",
            "server_test_darwin.go",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "failfast_linux.go",
            "server_test_linux.go",
        ],
        "@io_bazel_rules_go//go/platform:nacl": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "failfast_other.go",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "failfast_other.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "github.com/obeattie/tcp-failfast",
    visibility = ["//visibility:public"],
    deps = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "//vendor/github.com/songgao/water:go_default_library",
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//vendor/github.com/songgao/water:go_default_library",
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
