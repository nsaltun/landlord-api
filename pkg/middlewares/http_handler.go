package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type HttpContext struct {
	r   *http.Request
	w   http.ResponseWriter
	ctx context.Context
}

type CustomHandler func(*HttpContext)

type Middleware func(CustomHandler) CustomHandler

func MiddlewareRunner(handler CustomHandler, middlFns ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpCtx := &HttpContext{
			r:   r,
			w:   w,
			ctx: r.Context(),
		}

		final := handler
		for _, middlFn := range middlFns {
			final = middlFn(final)
		}

		final(httpCtx)
	}
}

type reqBodyCtxKey struct{}

var reqBodyKey = reqBodyCtxKey{}

func (c *HttpContext) ParseRequestBody() []byte {
	var bodyBytes []byte
	if c.r.Body != nil {
		defer c.r.Body.Close()
		var err error
		bodyBytes, err = io.ReadAll(c.r.Body)
		if err != nil {
			slog.Warn("error while parsing request body", "err", err)
			return nil
		}
		c.ctx = context.WithValue(c.ctx, reqBodyKey, bodyBytes)
	}

	return bodyBytes
}

func (c *HttpContext) CastRequestBody(v interface{}) error {
	bodyBytes := c.ctx.Value(reqBodyKey)
	if bodyBytes == nil {
		return errors.New("req body is nil")
	}
	err := json.Unmarshal(bodyBytes.([]byte), v)
	if err != nil {
		slog.Warn("request unmarshal error.", "err", err)
		return err
	}
	return nil
}

func (c *HttpContext) JSONErrResp(err error) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.Header().Set("statusCode", fmt.Sprint(http.StatusInternalServerError))
	c.w.WriteHeader(http.StatusInternalServerError)
	errResp := map[string]string{
		"errMsg": err.Error(),
	}
	if err := json.NewEncoder(c.w).Encode(errResp); err != nil {
		slog.Info("error while encoding the response", "err", err)
	}
}

// JSON sends a JSON response with the given status code
func (c *HttpContext) JSON(status int, payload interface{}) error {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.Header().Set("statusCode", fmt.Sprint(status))
	c.w.WriteHeader(status)
	json.NewEncoder(c.w).Encode(payload)
	return nil
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status
	}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	lrw.body = append(lrw.body, data...) // Capture response body
	return lrw.ResponseWriter.Write(data)
}

func LoggingMiddleware(next CustomHandler) CustomHandler {
	return func(c *HttpContext) {
		bodyBytes := c.ParseRequestBody()

		slog.Info("Request:",
			"method", c.r.Method,
			"path", c.r.URL.Path,
			"req", string(bodyBytes),
			"headers.ContentType", c.r.Header["Content-Type"],
			"headers.ContentLength", c.r.Header["Content-Length"],
			"headers.UserAgent", c.r.Header["User-Agent"],
		)

		lrw := newLoggingResponseWriter(c.w)
		c.w = lrw

		// Call the next handler
		next(c)

		// Log response details
		slog.Info("Response: ",
			"status", lrw.statusCode,
			"body", string(lrw.body))
	}
}
