package pkg

import (
	"log"

	"github.com/go-kit/kit/metrics/prometheus"
	kitLog "github.com/go-kit/log"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/repository"
	"github.com/ppeymann/vendors.git/server"
	"github.com/ppeymann/vendors.git/services/user"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func InitUserService(db *gorm.DB, logger kitLog.Logger, conf *config.Configuration, server *server.Server) models.UserService {
	repo := repository.NewUserRepo(db, "user")
	err := repo.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	// userService create service
	userService := user.NewService(repo, conf)

	// getting path
	path := getSchemaPath("user")
	userService, err = user.NewValidationService(userService, path)
	if err != nil {
		log.Fatal(err)
	}

	// @Injection logging service to chain
	userService = user.NewLoggerService(userService, kitLog.With(logger, "component", "user"))

	// @Injection Instrumenting service to chain
	userService = user.NewInstrumentingService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "user",
			Name:      "request_count",
			Help:      "number of request received. ",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "user",
			Name:      "request_latency_microseconds",
			Help:      "total duration of request (ms). ",
		}, fieldKeys),
		userService,
	)

	// @Injection Authorization service to chain
	userService = user.NewAuthService(userService)

	_ = user.NewHandler(userService, server)

	return userService
}
