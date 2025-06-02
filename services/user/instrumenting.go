package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/metrics"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type instrumentingService struct {
	next           models.UserService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, srv models.UserService) models.UserService {
	return &instrumentingService{
		next:           srv,
		requestCount:   requestCount,
		requestLatency: requestLatency,
	}
}

func (i *instrumentingService) Register(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "Register").Add(1)
		i.requestLatency.With("method", "Register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Register(ctx, in)
}
