package middleware

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status", "environment"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"method", "path", "status", "environment"})
)

var environment = detectEnvironment()

func detectEnvironment() string {
	ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return "local"
	}
	return string(ns)
}

// PrometheusMetrics tracks HTTP request count and latency per endpoint.
func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = "/unknown"
		}
		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status, environment).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path, status, environment).Observe(duration)
	}
}

type accessLogEntry struct {
	Level    string `json:"level"`
	Ts       float64 `json:"ts"`
	Msg      string `json:"msg"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Status   int    `json:"status"`
	Latency  string `json:"latency"`
	ClientIP string `json:"clientIP"`
	Service  string `json:"service"`
	Error    string `json:"error,omitempty"`  // NEW: error message if present
}

// AccessLogger writes one structured JSON line per request to stdout.
// Loki / Promtail will pick these up and index level, path, service as labels.
func AccessLogger() gin.HandlerFunc {
	encoder := json.NewEncoder(os.Stdout)
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		lvl := "info"
		msg := "request"
		errMsg := ""

		if c.Writer.Status() >= 400 {
			lvl = "warn"
			msg = "request-warning"
		}
		if c.Writer.Status() >= 500 {
			lvl = "error"
			msg = "request-error"
		}

		// Extract error details from gin.Context if present
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		encoder.Encode(accessLogEntry{
			Level:    lvl,
			Ts:       float64(time.Now().UnixNano()) / 1e9,
			Msg:      msg,
			Method:   c.Request.Method,
			Path:     path,
			Status:   c.Writer.Status(),
			Latency:  time.Since(start).String(),
			ClientIP: c.ClientIP(),
			Service:  "goldgym",
			Error:    errMsg,
		})
	}
}
