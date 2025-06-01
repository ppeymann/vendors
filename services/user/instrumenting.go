package user

import (
	"github.com/go-kit/kit/metrics"
	"github.com/ppeymann/vendors.git/models"
)

type instrumentingService struct {
	next           models.UserService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func NewInstrumentingService(srv models.UserService, requestCount metrics.Counter, requestLatency metrics.Histogram) models.UserService {
	return &instrumentingService{
		next:           srv,
		requestCount:   requestCount,
		requestLatency: requestLatency,
	}
}
