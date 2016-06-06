/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/06/06 10:01
 */

package wotracer

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/opentracing/basictracer-go"
	"github.com/opentracing/opentracing-go"

	"sourcegraph.com/sourcegraph/appdash"
	apptracer "sourcegraph.com/sourcegraph/appdash/opentracing"

	"github.com/wothing/wotracer/helper"
)

var tracer opentracing.Tracer

// InjectTracer MUST be run before all other func
func InitTracer(address string) {
	collector := appdash.NewRemoteCollector(address)
	tracer = apptracer.NewTracer(collector)
	opentracing.InitGlobalTracer(tracer)
}

// InjectRPC starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a parent. If no such parent could be found,
// StartSpanFromContext creates a root (parentless) Span.
func InjectRPC(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	return span, PackCtx(ctx)
}

// JoinRPC returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
func JoinRPC(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		// TODO log warn, but nothing we can do
	}
	ctx = helper.FromGRPCRequest(tracer, operationName)(ctx, &md)
	return opentracing.SpanFromContext(ctx), ctx
}

// PackCtx pack usual context to grpc context using metadata
func PackCtx(ctx context.Context) context.Context {
	md := metadata.Pairs()
	helper.ToGRPCRequest(tracer)(ctx, &md)

	return metadata.NewContext(
		ctx,
		md,
	)
}

// GetTraceID returns TraceID with given span
func GetTraceID(span opentracing.Span) string {
	return fmt.Sprintf("%016x", span.(basictracer.Span).Context().TraceID)
}
