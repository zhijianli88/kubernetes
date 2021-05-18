/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httptrace

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"k8s.io/klog/v2"

	apitrace "go.opentelemetry.io/otel/api/trace"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type contextKeyType int

const spanContextAnnotationKey string = "trace.kubernetes.io/context"

func stringToSpanContext(sc string) apitrace.SpanContext {
	id, _ := apitrace.IDFromHex(sc[0:32])
	spanid, _ := apitrace.SpanIDFromHex(sc[33:49])
	return apitrace.SpanContext{
		TraceID: id,
		SpanID:  spanid,
	}
}

// WithObject returns a context attached with a Span retrieved from object annotation, it doesn't start a new span
func WithObject(ctx context.Context, meta metav1.Object) context.Context {
	var latestContext string
	var latestTimeStamp *metav1.Time

	managedFields := meta.GetManagedFields()
	for _, mf := range managedFields {
		if latestTimeStamp != nil {
			if latestTimeStamp.Before(mf.Time) {
				latestTimeStamp = mf.Time
				latestContext = mf.TraceContext
			}
		} else {
			latestTimeStamp = mf.Time
			latestContext = mf.TraceContext
		}
		klog.V(3).InfoS("Trace request", "object", klog.KObj(meta), "Generation", meta.GetGeneration(), "trace-id", mf.TraceContext)
	}

	span := httpTraceSpan{
		spanContext: stringToSpanContext(latestContext),
	}
	klog.V(3).InfoS("Trace request", "object", klog.KObj(meta), "trace-id", latestContext)
	return apitrace.ContextWithSpan(ctx, span)
	// return spanContextFromAnnotations(ctx, meta, meta.GetAnnotations())
}

// spanContextFromAnnotations get span context from annotations
func spanContextFromAnnotations(ctx context.Context, meta metav1.Object, annotations map[string]string) context.Context {
	// get span context from annotations
	spanContext, err := decodeSpanContext(annotations[spanContextAnnotationKey])
	if err != nil {
		return ctx
	}
	span := httpTraceSpan{
		spanContext: spanContext,
	}
	klog.V(3).InfoS("Trace request", "object", klog.KObj(meta), "trace-id", spanContextString(spanContext))
	return apitrace.ContextWithSpan(ctx, span)
}

func spanContextString(spanContext apitrace.SpanContext) string {
	return fmt.Sprintf("%s-%s-%02d", spanContext.TraceID, spanContext.SpanID, spanContext.TraceFlags)
}

func StringSpanContextFromObject(meta metav1.Object) string {
	spanContext, err := decodeSpanContext(meta.GetAnnotations()[spanContextAnnotationKey])
	if err != nil {
		return ""
	}
	return spanContextString(spanContext)
}

// decodeSpanContext decode encodedSpanContext to spanContext
func decodeSpanContext(encodedSpanContext string) (apitrace.SpanContext, error) {
	// decode to byte
	byteList := make([]byte, base64.StdEncoding.DecodedLen(len(encodedSpanContext)))
	l, err := base64.StdEncoding.Decode(byteList, []byte(encodedSpanContext))
	if err != nil {
		return apitrace.EmptySpanContext(), err
	}
	byteList = byteList[:l]
	// decode to span context
	buffer := bytes.NewBuffer(byteList)
	spanContext := apitrace.SpanContext{}
	err = binary.Read(buffer, binary.LittleEndian, &spanContext)
	if err != nil {
		return apitrace.EmptySpanContext(), err
	}
	return spanContext, nil
}
