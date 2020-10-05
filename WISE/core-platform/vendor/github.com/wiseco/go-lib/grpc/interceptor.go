package grpc

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/wiseco/go-lib/log"
)

const (
	requestIDKey = ctxKey("requestID")
	loggerKey    = ctxKey("rLogger")
)

type ctxKey string

func unaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	var h interface{}

	start := time.Now()

	l := log.NewLogger()

	ri := getRequestID(l)

	defer func(ctx context.Context, l log.Logger, req interface{}, err error, ri string) {
		if r := recover(); r != nil {
			err = grpc.Errorf(codes.Internal, "panic: %v", r)

			f := log.Fields{
				"error":      err.Error(),
				"method":     info.FullMethod,
				"request_id": ri,
			}

			l.ErrorD("GRPC PANIC", f)
		}
	}(ctx, l, req, err, ri)

	ctx = context.WithValue(ctx, requestIDKey, ri)
	ctx = context.WithValue(ctx, loggerKey, l)

	h, err = handler(ctx, req)

	duration := time.Now().Sub(start)

	if err != nil {
		l.ErrorD("GRPC ERROR", log.Fields{
			"method":     info.FullMethod,
			"error":      err,
			"request_id": ri,
			"duration":   duration,
		})
	} else {
		l.InfoD("ACCESS", log.Fields{
			"method":     info.FullMethod,
			"request_id": ri,
			"duration":   duration,
		})
	}

	return h, err
}

func getRequestID(l log.Logger) string {
	hostname, err := os.Hostname()

	if err != nil || hostname == "" {
		hostname = "localhost"
	}

	u := uuid.New()

	return fmt.Sprintf("%s-%s", hostname, u.String())
}
