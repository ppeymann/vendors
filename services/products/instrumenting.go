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

func (i *instrumentingService) EditProduct(ctx *gin.Context, id uint, in *models.ProductInput) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCounter.With("method", "EditProduct").Add(1)
		i.requestLatency.With("method", "EditProduct").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.EditProduct(ctx, id, in)
}

func (i *instrumentingService) GetByTags(ctx *gin.Context, tags []string) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCounter.With("method", "GetByTags").Add(1)
		i.requestLatency.With("method", "GetByTags").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.GetByTags(ctx, tags)
}

func (i *instrumentingService) GetProduct(ctx *gin.Context, id uint) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCounter.With("method", "GetProduct").Add(1)
		i.requestLatency.With("method", "GetProduct").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.GetProduct(ctx, id)
}

// Add implements models.ProductService.
func (i *instrumentingService) Add(ctx *gin.Context, in *models.ProductInput) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCounter.With("method", "Add").Add(1)
		i.requestLatency.With("method", "Add").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Add(ctx, in)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, srv models.ProductService) models.ProductService {
	return &instrumentingService{
		requestCounter: counter,
		requestLatency: latency,
		next:           srv,
	}
}
