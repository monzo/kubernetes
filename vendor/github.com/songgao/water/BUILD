load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "if.go",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "params_darwin.go",
            "syscalls_darwin.go",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "if_linux.go",
            "params_linux.go",
            "syscalls_linux.go",
        ],
        "@io_bazel_rules_go//go/platform:nacl": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "params_others.go",
            "syscalls_other.go",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "params_windows.go",
            "syscalls_windows.go",
        ],
        "//conditions:default": [],
    }),
    importpath = "github.com/songgao/water",
    visibility = ["//visibility:public"],
    deps = select({
        "@io_bazel_rules_go//go/platform:windows": [
            "//vendor/golang.org/x/sys/windows/registry:go_default_library",
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
