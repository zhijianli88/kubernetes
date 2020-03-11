// // Package traceutil provides various definitions and utilities that allow for
// // common operations with our trace tooling, such as span creation, encoding, decoding,
// // and enumeration of possible services.
// package traceutil
package traceutil

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"contrib.go.opencensus.io/exporter/ocagent"

	flag "github.com/spf13/pflag"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// TraceAnnotationKey is the annotation name where span context should be found
const TraceAnnotationKey string = "trace.kubernetes.io/context"

type RequestID struct {
	enabled bool
}

var requestID RequestID

func Enabled() bool {
	return requestID.enabled
}

func InitFlags(flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.CommandLine
	}

	flagset.BoolVar(&requestID.enabled, "request-id", true, "If true, built-in request-id in log")
}

// InitializeExporter takes a ServiceType and sets the global OpenCensus exporter
// to export to that service on a specified Zipkin instance
func InitializeExporter(service string) {
	klog.Infof("OpenCensus trace exporter initializing with service %s", string(service))
	// create ocagent exporter
	exp, err := ocagent.NewExporter(ocagent.WithInsecure(), ocagent.WithServiceName(string(service)))
	if err != nil {
		log.Fatalf("Failed to create the agent exporter: %v", err)
	}
	// Only sample when the propagated parent SpanContext is sampled
	// Use ProbabilitySampler because it propagates the parent sampling decision.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0)})

	trace.RegisterExporter(exp)

	return
}

func ReqIDFromObj(ctx context.Context, tracedResource meta.Object) trace.TraceID {
	ctx, span := StartSpanFromObject(ctx, tracedResource, "object")
	defer span.End()
	EncodeContextIntoObject(ctx, tracedResource)

	return span.SpanContext().TraceID
}

func LogReqIDFromObject(ctx context.Context, tracedResource meta.Object, format string, args ...interface{}) {
	if requestID.enabled {
		var buffer bytes.Buffer
		reqid := ReqIDFromCtx(ctx, tracedResource)
		buffer.WriteString(fmt.Sprintf("RequestID[%s]: ", reqid))
		buffer.WriteString(fmt.Sprintf(format, args...))
		klog.Info(buffer.String())
	} else {
		klog.Infof(format, args...)
	}
}

func ReqIDFromCtx(ctx context.Context, tracedResource meta.Object) trace.TraceID {
	c, span := trace.StartSpan(ctx, "name", trace.WithSampler(trace.AlwaysSample()))
	defer span.End()
	EncodeContextIntoObject(c, tracedResource)

	return span.SpanContext().TraceID
}

func LogReqIDFromContext(ctx context.Context, tracedResource meta.Object, format string, args ...interface{}) {
	if requestID.enabled {
		var buffer bytes.Buffer
		reqid := ReqIDFromCtx(ctx, tracedResource)
		buffer.WriteString(fmt.Sprintf("RequestID[%s]: ", reqid))
		buffer.WriteString(fmt.Sprintf(format, args...))
		klog.Info(buffer.String())
	} else {
		klog.Infof(format, args...)
	}
}

// StartSpanFromObject takes an object to extract trace context from and the desired Span name and
// constructs a new Span from this information.  It mirrors trace.StartSpan, but for kubernetes objects.
func StartSpanFromObject(ctx context.Context, tracedResource meta.Object, name string) (context.Context, *trace.Span) {
	spanFromEncodedContext, ok := spanContextFromObject(tracedResource)
	if !ok {
		return ctx, &trace.Span{}
	}
	return trace.StartSpanWithRemoteParent(ctx, name, spanFromEncodedContext)
}

// spanContextFromObject takes an object to extract an encoded SpanContext from and returns the decoded SpanContext
func spanContextFromObject(tracedResource meta.Object) (trace.SpanContext, bool) {
	tracedResourceAnnotations := tracedResource.GetAnnotations()
	embeddedSpanContext, ok := tracedResourceAnnotations[TraceAnnotationKey]
	if !ok {
		return trace.SpanContext{}, false
	}

	decodedContextBytes, err := base64.StdEncoding.DecodeString(embeddedSpanContext)
	if err != nil {
		return trace.SpanContext{}, false
	}

	return propagation.FromBinary(decodedContextBytes)

}

// EncodeContextIntoObject encodes the SpanContext contained in the context into the provided object
func EncodeContextIntoObject(ctx context.Context, tracedResource meta.Object) {
	klog.Infof("OC trace:EncodeContextIntoObject")
	span := trace.FromContext(ctx)
	if span != nil {
		encodeSpanContextIntoObject(span.SpanContext(), tracedResource)
		klog.Infof("OC trace:EncodeContextIntoObject : TraceID:%s", span.SpanContext().TraceID)
		tracedResourceAnnotations := tracedResource.GetAnnotations()
		klog.Infof("OC trace:EncodeContextIntoObject : Annotation?: %s", tracedResourceAnnotations[TraceAnnotationKey])
	}
}

func RemoveSpanContextFromObject(tracedResource meta.Object) {
	klog.Infof("OC trace:RemoveSpanContextFromObject")

	tracedResourceAnnotations := tracedResource.GetAnnotations()
	klog.Infof("OC trace:RemoveSpanContextFromObject : Annotation?: %s", tracedResourceAnnotations[TraceAnnotationKey])
	delete(tracedResourceAnnotations, TraceAnnotationKey)
	tracedResource.SetAnnotations(tracedResourceAnnotations)
}

// encodeSpanContextIntoObject takes a pointer to an object and a Span Context to embed
// Base64 encodes the wire format for the SpanContext, and puts it in the object's TraceContext field
func encodeSpanContextIntoObject(ctx trace.SpanContext, tracedResource meta.Object) {
	tracedResourceAnnotations := tracedResource.GetAnnotations()

	rawContextBytes := propagation.Binary(ctx)
	encodedContext := base64.StdEncoding.EncodeToString(rawContextBytes)

	tracedResourceAnnotations[TraceAnnotationKey] = encodedContext
	tracedResource.SetAnnotations(tracedResourceAnnotations)

	return
}

// import (
// 	"context"
// 	"encoding/base64"
// 	"log"

// 	// "contrib.go.opencensus.io/exporter/ocagent"

// 	// "go.opencensus.io/trace"
// 	// "go.opencensus.io/trace/propagation"

// 	"k8s.io/klog"
// 	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

// 	"go.opentelemetry.io/otel/api/core"
// 	"go.opentelemetry.io/otel/api/global"
// 	"go.opentelemetry.io/otel/api/propagators"
// 	"go.opentelemetry.io/otel/api/trace"

// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"

// 	"go.opentelemetry.io/otel/exporter/trace/stdout"
// )

// // TraceAnnotationKey is the annotation name where span context should be found
// const TraceAnnotationKey string = "trace.kubernetes.io/context"

// // InitializeExporter takes a ServiceType and sets the global OpenCensus exporter
// // to export to that service on a specified Zipkin instance
// func InitializeExporter(service string) {
// 	klog.Infof("OpenCensus trace exporter initializing with service %s", string(service))

// 	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
// 		sdktrace.WithSyncer(exporter))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	global.SetTraceProvider(tp)
// }

// // StartSpanFromObject takes an object to extract trace context from and the desired Span name and
// // constructs a new Span from this information.  It mirrors trace.StartSpan, but for kubernetes objects.
// func StartSpanFromObject(ctx context.Context, tracedResource meta.Object, name string) (context.Context, *sdktrace.Span) {
// 	spanFromEncodedContext, ok := spanContextFromObject(tracedResource)
// 	if !ok {
// 		return ctx, &sdktrace.Span{}
// 	}
// 	tr := global.TraceProvider().Tracer("trace/traceutil")
// 	//return trace.StartSpanWithRemoteParent(ctx, name, spanFromEncodedContext)
// 	return tr.Start(
// 		trace.ContextWithRemoteSpanContext(ctx, spanFromEncodedContext),name)
// }

// // spanContextFromObject takes an object to extract an encoded SpanContext from and returns the decoded SpanContext
// func spanContextFromObject(tracedResource meta.Object) (core.SpanContext, bool) {
// 	tracedResourceAnnotations := tracedResource.GetAnnotations()
// 	embeddedSpanContext, ok := tracedResourceAnnotations[TraceAnnotationKey]
// 	if !ok {
// 		return core.SpanContext{}, false
// 	}

// 	decodedContextBytes, err := base64.StdEncoding.DecodeString(embeddedSpanContext)
// 	if err != nil {
// 		return core.SpanContext{}, false
// 	}

// 	//return propagation.FromBinary(decodedContextBytes)
// 	return propagators.FromBytes(decodedContextBytes)

// }

// // EncodeContextIntoObject encodes the SpanContext contained in the context into the provided object
// func EncodeContextIntoObject(ctx context.Context, tracedResource meta.Object) {
// 	span := sdktrace.FromContext(ctx)
// 	if span != nil {
// 		encodeSpanContextIntoObject(span.SpanContext(), tracedResource)
// 	}
// }

// func RemoveSpanContextFromObject(tracedResource meta.Object) {
// 	tracedResourceAnnotations := tracedResource.GetAnnotations()
// 	delete(tracedResourceAnnotations, TraceAnnotationKey)
// 	tracedResource.SetAnnotations(tracedResourceAnnotations)
// }

// // encodeSpanContextIntoObject takes a pointer to an object and a Span Context to embed
// // Base64 encodes the wire format for the SpanContext, and puts it in the object's TraceContext field
// func encodeSpanContextIntoObject(ctx core.SpanContext, tracedResource meta.Object) {
// 	tracedResourceAnnotations := tracedResource.GetAnnotations()

// 	rawContextBytes := propagators.Binary(ctx)
// 	encodedContext := base64.StdEncoding.EncodeToString(rawContextBytes)

// 	tracedResourceAnnotations[TraceAnnotationKey] = encodedContext
// 	tracedResource.SetAnnotations(tracedResourceAnnotations)

// 	return
// }
