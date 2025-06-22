package pkg

import (
	"log"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitLog "github.com/go-kit/log"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/repository"
	"github.com/ppeymann/vendors.git/server"
	"github.com/ppeymann/vendors.git/services/mio"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func InitMioService(db *gorm.DB, logger kitLog.Logger, conf *config.Configuration, srv *server.Server) models.MioService {
	// Create repository
	mioRepo, err := repository.NewMioRepo(db, conf.Storage, env.GetEnv("DATABASE", ""))
	if err != nil {
		log.Fatal(err)
	}

	// Migrate this table
	err = mioRepo.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	// Create Service
	mioService := mio.NewService(conf.Storage, mioRepo)

	// get path and inject it to validation
	path := getSchemaPath("mio")
	mioService, err = mio.NewValidationService(conf.Storage, mioService, env.GetEnv("JWT", ""), path)
	if err != nil {
		log.Fatal(err)
	}

	// @Inject logging service to chain
	mioService = mio.NewLoggerService(kitLog.With(logger, "component", "mio"), mioService)

	// @Inject instrumenting service to chain
	mioService = mio.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "mio",
			Name:      "request_count",
			Help:      "num of requests received",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "mio",
			Name:      "request_latency_microseconds",
			Help:      "total duration of request (ms).",
		}, fieldKeys),
		mioService,
	)

	// @Inject authorization service to chain
	mioService = mio.NewAuthService(conf.Storage, mioService)

	// Create api handler
	_ = mio.NewHandler(mioService, conf, srv)

	// return service
	return mioService
}
