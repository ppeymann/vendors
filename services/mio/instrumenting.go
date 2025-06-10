package mio

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/metrics"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           models.MioService
}

// Download implements models.MioService.
func (i *instrumentingService) Download(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "download").Add(1)
		i.requestLatency.With("method", "download").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Download(in, ctx)
}

// Image implements models.MioService.
func (i *instrumentingService) Image(in *models.DownloadInput, ctx *gin.Context) ([]byte, *models.StorageEntity, error) {
	defer func(begin time.Time) {
		i.requestCount.With("method", "image").Add(1)
		i.requestLatency.With("method", "image").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Image(in, ctx)
}

// Upload implements models.MioService.
func (i *instrumentingService) Upload(in *models.UploadInput, ctx *gin.Context) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "Upload").Add(1)
		i.requestLatency.With("method", "Upload").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Upload(in, ctx)
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, service models.MioService) models.MioService {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           service,
	}
}
