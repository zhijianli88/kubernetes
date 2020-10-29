load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "attribute.go",
        "instrumentation.go",
        "metric.go",
        "resource.go",
        "span.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/go.opentelemetry.io/otel/exporters/otlp/internal/transform",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/internal/transform",
    visibility = ["//vendor/go.opentelemetry.io/otel/exporters/otlp:__subpackages__"],
    deps = [
        "//vendor/go.opentelemetry.io/otel/api/metric:go_default_library",
        "//vendor/go.opentelemetry.io/otel/api/trace:go_default_library",
        "//vendor/go.opentelemetry.io/otel/codes:go_default_library",
        "//vendor/go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/common/v1:go_default_library",
        "//vendor/go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/metrics/v1:go_default_library",
        "//vendor/go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/resource/v1:go_default_library",
        "//vendor/go.opentelemetry.io/otel/exporters/otlp/internal/opentelemetry-proto-gen/trace/v1:go_default_library",
        "//vendor/go.opentelemetry.io/otel/label:go_default_library",
        "//vendor/go.opentelemetry.io/otel/sdk/export/metric:go_default_library",
        "//vendor/go.opentelemetry.io/otel/sdk/export/metric/aggregation:go_default_library",
        "//vendor/go.opentelemetry.io/otel/sdk/export/trace:go_default_library",
        "//vendor/go.opentelemetry.io/otel/sdk/instrumentation:go_default_library",
        "//vendor/go.opentelemetry.io/otel/sdk/resource:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
