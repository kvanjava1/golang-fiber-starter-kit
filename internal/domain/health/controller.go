package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"

	"golang-starter-kit/pkg/database"
	"golang-starter-kit/pkg/logger"
	"golang-starter-kit/pkg/response"
)

type Controller struct {
	dbMgr       *database.DBManager
	redisClient *redis.Client
}

func NewController(dbMgr *database.DBManager, redisClient *redis.Client) *Controller {
	return &Controller{
		dbMgr:       dbMgr,
		redisClient: redisClient,
	}
}

type ComponentStatus struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

type HealthResponse struct {
	Status     string                     `json:"status"`
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentStatus `json:"components"`
}

// Check memeriksa kesehatan koneksi database dan redis
func (ctrl *Controller) Check(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	overallStatus := "healthy"
	components := make(map[string]ComponentStatus)

	// 1. Check SQL Database
	if ctrl.dbMgr != nil {
		sqlDB := ctrl.dbMgr.DefaultSQL()
		if sqlDB != nil {
			start := time.Now()
			err := sqlDB.PingContext(ctx)
			latency := time.Since(start).String()

			if err != nil {
				overallStatus = "unhealthy"
				components["database"] = ComponentStatus{
					Status: "down",
					Error:  err.Error(),
				}
			} else {
				components["database"] = ComponentStatus{
					Status:  "up",
					Latency: latency,
				}
			}
		}
	}

	// 2. Check Redis (jika dikonfigurasi)
	if ctrl.redisClient != nil {
		start := time.Now()
		err := ctrl.redisClient.Ping(ctx).Err()
		latency := time.Since(start).String()

		if err != nil {
			overallStatus = "degraded"
			components["redis"] = ComponentStatus{
				Status: "down",
				Error:  err.Error(),
			}
		} else {
			components["redis"] = ComponentStatus{
				Status:  "up",
				Latency: latency,
			}
		}
	}

	traceID := logger.CorrelationIDFromContext(c.Context())
	meta := response.NewMeta(traceID)

	res := HealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Components: components,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(response.Success(overallStatus, res, meta))
}
