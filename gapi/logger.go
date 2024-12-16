package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err) // this allows us to add more information to the log
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode                // we are adding to the current implementation
	rec.ResponseWriter.WriteHeader(statusCode) // since we are adding to the original code, we mean call the origin func
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		/*
			no result or err object from the http.ServeHTTP, so we need a custom implementation
			of the interface to add code to track errors and results
		*/
		recorder := &ResponseRecorder{
			ResponseWriter: response,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(recorder, request)
		duration := time.Since(startTime)

		logger := log.Info()
		if recorder.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", recorder.Body)
		}

		logger.Str("protocol", "http").
			Str("method", request.Method).
			Str("path", request.RequestURI).
			Int("status_code", recorder.StatusCode).
			Str("status_text", http.StatusText(recorder.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
