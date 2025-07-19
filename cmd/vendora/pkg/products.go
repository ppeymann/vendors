package pkg

import (
	"github.com/go-kit/kit/metrics/prometheus"
	kitLog "github.com/go-kit/log"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/repository"
	"github.com/ppeymann/vendors.git/server"
	"github.com/ppeymann/vendors.git/services/products"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"log"
)

func InitProducts(db *gorm.DB, logger kitLog.Logger, conf *config.Configuration, server *server.Server) models.ProductService {
	repo := repository.NewProductsRepo(db, env.GetEnv("POSTGRES_DB", ""))
	err := repo.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	// productsService create service
	productsService := products.NewService(repo)

	// getting path
	path := getSchemaPath("products")
	productsService, err = products.NewValidationsService(path, productsService)
	if err != nil {
		log.Fatal(err)
	}

	// @Injection logger service to chain
	productsService = products.NewLoggerService(logger, productsService)

	// @Injection Instrumenting Service to chain
	productsService = products.NewInstrumentingService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "products",
			Name:      "request_count",
			Help:      "number of request received",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "products",
			Name:      "request_latency_microseconds",
			Help:      "total duration of request (ms). ",
		}, fieldKeys),
		productsService,
	)

	// @Injection Authorization service to chain
	productsService = products.NewAuthService(productsService)

	_ = products.NewHandler(productsService, server)

	return productsService
}
