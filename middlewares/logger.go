package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"

	"go-echo-api/zplogger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggerConfig configures the logging middleware.
type LoggerConfig struct {
	Logger         *zplogger.Logger
	IgnorePaths    []string
	IgnoreBodyKeys []string
	MaskedKeys     []string
}

type responseBodyWriter struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware logs every request with request ID, body, and response body.
func LoggerMiddleware(cfg LoggerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Response().Header().Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = uuid.New().String()
			}
			c.Response().Header().Set(echo.HeaderXRequestID, rid)

			req := c.Request()
			path := req.URL.Path
			method := req.Method
			uri := req.RequestURI
			ip := c.RealIP()

			// Skip ignored paths entirely
			for _, p := range cfg.IgnorePaths {
				if strings.HasPrefix(path, p) || path == p {
					return next(c)
				}
			}

			// Read request body (skip multipart)
			var reqBodyStr string
			mediaType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if mediaType == "multipart/form-data" {
				reqBodyStr = "[multipart/form-data]"
			} else if req.Body != nil {
				bodyBytes, _ := io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				if len(bodyBytes) > 0 {
					reqBodyStr = sanitizeBody(bodyBytes, cfg.IgnoreBodyKeys, cfg.MaskedKeys)
				}
			}

			// Wrap response writer
			w := &responseBodyWriter{ResponseWriter: c.Response().Writer}
			c.Response().Writer = w

			fields := []zap.Field{
				zap.String("requestId", rid),
				zap.String("method", method),
				zap.String("uri", uri),
				zap.String("ip", ip),
			}
			if reqBodyStr != "" {
				fields = append(fields, zap.String("reqbody", reqBodyStr))
			}
			cfg.Logger.Info("Request", fields...)

			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			respFields := []zap.Field{
				zap.String("requestId", rid),
				zap.Int("status", c.Response().Status),
				zap.Int64("latency_ms", latency.Milliseconds()),
			}

			// Response body logging
			respBody := w.body.Bytes()
			if len(respBody) > 0 {
				if len(respBody) < 1024 {
					sanitized := sanitizeBody(respBody, cfg.IgnoreBodyKeys, cfg.MaskedKeys)
					respFields = append(respFields, zap.String("resbody", sanitized))
				} else {
					respFields = append(respFields, zap.Int("resByteLength", len(respBody)))
				}
			}

			if err != nil {
				respFields = append(respFields, zap.Error(err))
			}

			// Combine into a single log line via zap.Any or append
			cfg.Logger.Info("Response", respFields...)

			return err
		}
	}
}

// sanitizeBody parses JSON body and removes/masks the specified keys.
func sanitizeBody(body []byte, ignoreBodyKeys, maskedKeys []string) string {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error sanitizeBody parses JSON body: %s", err.Error())
		return string(body)
	}

	for _, key := range ignoreBodyKeys {
		delete(data, key)
	}
	for _, key := range maskedKeys {
		if _, ok := data[key]; ok {
			data[key] = "***"
		}
	}

	out, _ := json.Marshal(data)
	return string(out)
}
