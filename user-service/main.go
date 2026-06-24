package main

import (
	"context"
	"log"
	"time"
	"strconv"

	"user-service/src/config/database"
	"user-service/src/config/jwt"
	handler "user-service/src/handler/user"
	"user-service/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total jumlah HTTP request",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Durasi HTTP request dalam seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func main() {
	if err := initDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	if err := runServerWithMetrics(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func runServerWithMetrics() error {
	defer func() {
		if err := database.CloseDatabase(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	cfg := jwt.LoadConfigJWT()
	authMiddleware := middleware.AuthMiddleware(middleware.Options{
		JWTSecret: cfg.JWTSecret,
	})

	r := gin.Default()

	r.Use(prometheusMiddleware())

	r.GET("/public", handler.PublicHandler)
	r.GET("/users/all", authMiddleware, handler.GetAllUserHandler)
	r.GET("/users/profile", authMiddleware, handler.GetUserProfileHandler)
	r.GET("/users/profile/:id", authMiddleware, handler.GetUserByIdHandler)
	r.PUT("/users/:id", authMiddleware, handler.UpdateUserHandler)
	r.DELETE("/users/:id", authMiddleware, handler.DeleteUserHandler)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Println("User Service is running on :3001")
	log.Println("Metrics available at http://localhost:3001/metrics")
	
	return r.Run(":3001")
}

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		endpoint := c.FullPath() 
		
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			endpoint,
			status,
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			endpoint,
		).Observe(duration)
	}
}

func initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return database.DatabaseConfigWithRetryContext(ctx, 3)
}