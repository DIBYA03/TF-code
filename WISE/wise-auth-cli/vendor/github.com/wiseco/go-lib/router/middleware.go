package router

import (
	"context"
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/middleware"

	"github.com/wiseco/go-lib/log"
)

const loggerKey = ctxKey("rLogger")

type ctxKey string

//Middleware is an interface describing middleware methods
type Middleware interface {
	Logger(http.Handler) http.Handler
}

type middleware struct{}

//NewMiddleWare returns a middleware interface
func NewMiddleWare() Middleware {
	return &middleware{}
}

//Logger is a middleware used for all requests, see NewRouter
func (m middleware) Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := log.NewLogger()

		r = addLoggerContextToRequest(l, r)

		sw := statusWriter{ResponseWriter: w}

		next.ServeHTTP(&sw, r)

		duration := time.Now().Sub(start)

		l.InfoD("ACCESS", log.Fields{
			"host":           r.Host,
			"method":         r.Method,
			"path":           r.URL.Path,
			"status":         sw.status,
			"request_id":     chiMiddleware.GetReqID(r.Context()),
			"content_length": sw.length,
			"duration":       duration,
		})
	}

	return http.HandlerFunc(fn)
}

//GetLogger returns the logger in the request context
func GetLogger(r *http.Request) log.Logger {
	ctx := r.Context()

	l, ok := ctx.Value(loggerKey).(log.Logger)

	//If the logger isn't present in the request(this should never happen)
	//Let's add a logger to the request, but log a warning
	if !ok {
		l = log.NewLogger()

		*r = *addLoggerContextToRequest(l, r)

		l.WarnD("Logger missing from request context", log.Fields{"path": r.URL.Path})
	}

	return l
}

func addLoggerContextToRequest(l log.Logger, r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, loggerKey, l)

	return r.WithContext(ctx)
}
