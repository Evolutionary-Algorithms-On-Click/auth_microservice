package util

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var SharedLogger *LoggerService

type LoggerService struct {
	Logger zerolog.Logger
	Env    string
}

// environments
const (
	EnvDevelopment = "DEVELOPMENT"
	EnvProduction  = "PRODUCTION"
)

type ILoggerService interface {
	// Function to enrich each log with data
	enrich(req *http.Request, e *zerolog.Event) *zerolog.Event

	// This set of functions is to be used in the context of the web-server
	// where there is a server context involved
	DebugCtx(req *http.Request, msg string)
	InfoCtx(req *http.Request, msg string)
	WarnCtx(req *http.Request, msg string)
	ErrorCtx(req *http.Request, msg string, err error)
	FatalCtx(req *http.Request, msg string, err error)
	PanicCtx(req *http.Request, msg string, r any, trace string) // r = recover()
	SuccessCtx(req *http.Request)

	// This set of functions can be used in scenarios where there is no
	// server context involved
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Fatal(msg string, err error)

	// Logging middleware to be used only as a global middleware during router
	// initialization
	LogMiddleware(next http.Handler) http.Handler
}

func InitLogger(env string) (*LoggerService, error) {
	if env == "" {
		env = EnvDevelopment
	}
	if env != EnvDevelopment && env != EnvProduction {
		return nil, fmt.Errorf("invalid environment for logger setup: %s", env)
	}

	config := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(config).With().Timestamp().Logger()
	return &LoggerService{
		Logger: logger,
		Env:    env,
	}, nil
}

func (l *LoggerService) enrich(req *http.Request, e *zerolog.Event) *zerolog.Event {
	queryParams := req.URL.Query()
	clientIP := req.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = req.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP, _, _ = net.SplitHostPort(req.RemoteAddr)
	}
	// In case of multiple IPs in X-Forwarded-For, take the first one
	if strings.Contains(clientIP, ",") {
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	}

	return e.
		Str("route", req.URL.Path).
		Str("method", req.Method).
		Interface("query-params", queryParams).
		Str("ip", clientIP).
		Str("user-agent", req.Header.Get("User-Agent"))
}

func (l *LoggerService) DebugCtx(req *http.Request, msg string) {
	if l.Env == EnvProduction {
		return
	}
	event := l.Logger.WithLevel(zerolog.DebugLevel)
	l.enrich(req, event).Msg(msg)
}

func (l *LoggerService) InfoCtx(req *http.Request, msg string) {
	event := l.Logger.WithLevel(zerolog.InfoLevel)
	l.enrich(req, event).Msg(msg)
}

func (l *LoggerService) WarnCtx(req *http.Request, msg string) {
	event := l.Logger.WithLevel(zerolog.WarnLevel)
	l.enrich(req, event).Msg(msg)
}

func (l *LoggerService) ErrorCtx(req *http.Request, msg string, err error) {
	event := l.Logger.WithLevel(zerolog.ErrorLevel).Err(err)
	l.enrich(req, event).Msg(msg)
}

func (l *LoggerService) FatalCtx(req *http.Request, msg string, err error) {
	event := l.Logger.WithLevel(zerolog.FatalLevel).Err(err)
	l.enrich(req, event).Msg(msg)
}

func (l *LoggerService) PanicCtx(req *http.Request, msg string, r any, trace string) {
	event := l.Logger.WithLevel(zerolog.InfoLevel).
		Str("panic_value", fmt.Sprintf("%v", r)).
		Str("trace", trace)
	l.enrich(req, event).Msg(msg)
}

func (l *LoggerService) SuccessCtx(req *http.Request) {
	event := l.Logger.WithLevel(zerolog.InfoLevel)
	l.enrich(req, event).Msg("request successful")
}

func (l *LoggerService) Debug(msg string) {
	if l.Env == EnvProduction {
		return
	}
	l.Logger.WithLevel(zerolog.DebugLevel).Msg(msg)
}

func (l *LoggerService) Info(msg string) {
	l.Logger.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *LoggerService) Warn(msg string) {
	l.Logger.WithLevel(zerolog.InfoLevel).Msg(msg)
}

func (l *LoggerService) Error(msg string, err error) {
	l.Logger.WithLevel(zerolog.InfoLevel).Err(err).Msg(msg)
}

func (l *LoggerService) Fatal(msg string, err error) {
	l.Logger.WithLevel(zerolog.FatalLevel).Err(err).Msg(msg)
}

// loggingResponseWriter is a wrapper around http.ResponseWriter to capture status code and response size.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

func (l *LoggerService) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 0, size: 0}
		next.ServeHTTP(lrw, r)

		// Get client IP from headers or remote address
		clientIP := r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.Header.Get("X-Real-IP")
		}
		if clientIP == "" {
			clientIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		// In case of multiple IPs in X-Forwarded-For, take the first one
		if strings.Contains(clientIP, ",") {
			clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		}

		event := l.Logger.WithLevel(zerolog.InfoLevel).
			Str("route", r.URL.Path).
			Str("method", r.Method).
			Int("status", lrw.statusCode).
			Int("response-size", lrw.size).
			Dur("duration", time.Since(start)).
			Interface("query-params", r.URL.Query()).
			Str("ip", clientIP).
			Str("user-agent", r.Header.Get("User-Agent"))

		event.Send()
	})
}

// Removed global SharedLogger to prevent duplicate logging
