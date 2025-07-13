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

// Login implements models.UserService.
func (i *instrumentingService) Login(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "Login").Add(1)
		i.requestLatency.With("method", "Login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Login(ctx, in)
}

// User implements models.UserService.
func (i *instrumentingService) User(ctx *gin.Context) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "User").Add(1)
		i.requestLatency.With("method", "User").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.User(ctx)
}

// EditUser implements models.UserService.
func (i *instrumentingService) EditUser(ctx *gin.Context, in *models.EditUserInput) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "EditUser").Add(1)
		i.requestLatency.With("method", "EditUser").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.EditUser(ctx, in)
}

// GetAllUserWithRole implements models.UserService.
func (i *instrumentingService) GetAllUserWithRole(ctx *gin.Context, role string) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetAllUserWithRole").Add(1)
		i.requestLatency.With("method", "GetAllUserWithRole").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.GetAllUserWithRole(ctx, role)
}

// ActiveDeActiveSuspended implements models.UserService.
func (i *instrumentingService) ActiveDeActiveSuspended(ctx *gin.Context) *vendora.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "ActiveDeActiveSuspended").Add(1)
		i.requestLatency.With("method", "ActiveDeActiveSuspended").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.ActiveDeActiveSuspended(ctx)
}
