package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	LocalMetadataKey = "metadata"
)

func HeaderFromIncoming(ctx context.Context, key string) []string {
	s := []string{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return s
		}
	}

	return md[key]
}

func HeaderValueFromIncoming(ctx context.Context, key string) string {
	var s string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return s
		}
	}

	val := md[key]
	if len(val) > 0 {
		return val[0]
	}

	return s
}

func AttachMetadataToNewOutgoingContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return ctx
		}
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func HeaderFromOutgoing(ctx context.Context, key string) []string {
	s := []string{}
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return s
		}
	}

	return md[key]
}

func HeaderValueFromOutgoing(ctx context.Context, key string) string {
	var s string
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return s
		}
	}

	val := md[key]
	if len(val) > 0 {
		return val[0]
	}

	return s
}

func MetadataFromIncoming(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return metadata.New(map[string]string{})
		}
	}

	return md
}

func MetadataFromOutgoing(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md, ok = ctx.Value(LocalMetadataKey).(metadata.MD)
		if !ok || md == nil {
			return metadata.New(map[string]string{})
		}
	}

	return md
}
