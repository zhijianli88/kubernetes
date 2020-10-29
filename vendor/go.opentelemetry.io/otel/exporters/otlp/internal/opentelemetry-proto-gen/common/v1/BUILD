load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["common.pb.go"],
    importmap = "k8s.io/kubernetes/vendor/go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/common/v1",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/common/v1",
    visibility = ["//vendor/go.opentelemetry.io/otel/exporters/otlp:__subpackages__"],
    deps = ["//vendor/github.com/gogo/protobuf/proto:go_default_library"],
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
