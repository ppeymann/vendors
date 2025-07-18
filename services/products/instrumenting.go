package products

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/metrics"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type instrumentingService struct {
	requestCounter metrics.Counter
	requestLatency metrics.Histogram
	next           models.ProductService
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, srv models.ProductService) models.ProductService {
	return &instrumentingService{
		requestCounter: counter,
		requestLatency: latency,
		next:           srv,
	}
}

// Add implements models.ProductService.
func (i *instrumentingService) Add(ctx *gin.Context, in *models.ProductInput) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCounter.With("method", "Add").Add(1)
		i.requestLatency.With("method", "Add").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Add(ctx, in)
}
