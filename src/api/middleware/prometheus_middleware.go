package middleware

import (
	utils "bookem-user-service/util"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status", "endpoint"},
	)

	httpResponseSizeBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_size_bytes",
			Help: "Total response size in bytes",
		},
		[]string{"endpoint", "status"},
	)
)

func PrometheusMiddleware() gin.HandlerFunc {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpResponseSizeBytes)

	return func(c *gin.Context) {
		c.Next()

		endpoint := c.FullPath()
		status := fmt.Sprintf("%d", c.Writer.Status())
		method := c.Request.Method
		size := float64(c.Writer.Size())

		httpRequestsTotal.WithLabelValues(method, status, endpoint).Inc()

		if size >= 0 {
			httpResponseSizeBytes.WithLabelValues(endpoint, status).Add(float64(size))
		} else {
			utils.TEL.Warn("Response size < 0, cannot push to Prometheus", "size", size)
		}
	}
}
