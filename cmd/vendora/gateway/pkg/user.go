package pkg

import (
	"fmt"
	"log"

	"github.com/go-kit/kit/metrics/prometheus"
	kitLog "github.com/go-kit/log"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/models"
	userpb "github.com/ppeymann/vendors.git/proto/user"
	"github.com/ppeymann/vendors.git/server"
	"github.com/ppeymann/vendors.git/services/user"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserService(logger kitLog.Logger, conf *config.Configuration, server *server.Server) models.UserService {
	connString := fmt.Sprintf("%s%s", env.GetEnv("GRPC_ADRR", "localhost"), env.GetEnv("USER_PORT", "50051"))

	// connection to gRPC Server
	conn, err := grpc.NewClient(connString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot connect to user-service gRPC: ", err)
	}

	grpcClient := userpb.NewUserServiceClient(conn)

	// userService create service
	userService := user.NewService(grpcClient, conf)

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

	// @Injection
	userService = user.NewAuthService(userService)

	_ = user.NewHandler(userService, server)

	return userService
}
