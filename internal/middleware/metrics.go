package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

//! \fn MetricsMiddleware() gin.HandlerFunc
//! \brief Tracks HTTP request metrics for Prometheus.
//! \return Gin middleware function.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.Request.URL.Path).Inc()
		c.Next()
	}
}

//! \var httpRequestsTotal
//! \brief Prometheus counter for total HTTP requests.
var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "endpoint"},
)

//! \var tasksCreatedTotal
//! \brief Prometheus counter for total tasks created.
var tasksCreatedTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "tasks_created_total",
		Help: "Total number of tasks created",
	},
)

//! \fn init()
//! \brief Registers Prometheus metrics.
func init() {
	prometheus.MustRegister(httpRequestsTotal, tasksCreatedTotal)
}