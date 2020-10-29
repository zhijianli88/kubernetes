load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "baggage.go",
        "doc.go",
        "otel.go",
        "propagation.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/go.opentelemetry.io/otel",
    importpath = "go.opentelemetry.io/otel",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/go.opentelemetry.io/otel/api/metric:go_default_library",
        "//vendor/go.opentelemetry.io/otel/api/trace:go_default_library",
        "//vendor/go.opentelemetry.io/otel/internal/baggage:go_default_library",
        "//vendor/go.opentelemetry.io/otel/label:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//vendor/go.opentelemetry.io/otel/api/global:all-srcs",
        "//vendor/go.opentelemetry.io/otel/api/metric:all-srcs",
        "//vendor/go.opentelemetry.io/otel/api/trace:all-srcs",
        "//vendor/go.opentelemetry.io/otel/codes:all-srcs",
        "//vendor/go.opentelemetry.io/otel/exporters/otlp:all-srcs",
        "//vendor/go.opentelemetry.io/otel/internal:all-srcs",
        "//vendor/go.opentelemetry.io/otel/label:all-srcs",
        "//vendor/go.opentelemetry.io/otel/propagators:all-srcs",
        "//vendor/go.opentelemetry.io/otel/sdk:all-srcs",
        "//vendor/go.opentelemetry.io/otel/semconv:all-srcs",
        "//vendor/go.opentelemetry.io/otel/unit:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
