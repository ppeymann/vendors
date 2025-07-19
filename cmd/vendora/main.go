package main

import (
	"fmt"
	"log"
	"os"
	"time"

	kitLog "github.com/go-kit/log"
	"github.com/ppeymann/vendors.git/cmd/vendora/pkg"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/server"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	now := time.Now().UTC()

	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Unix()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 35, 0, 0, time.UTC).Unix()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 30, 0, 0, time.UTC).Unix()

	fmt.Println("date:", base, "start:", start, "end:", end)

	config, err := config.NewConfiguration(env.GetEnv("SESSION", "session_secret_string__configs"))
	if err != nil {
		log.Fatal(err)

		return
	}

	db, err := gorm.Open(pg.Open(env.GetEnv("DSN", "")), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)

		return
	}

	// configuration logger
	var logger kitLog.Logger
	logger = kitLog.NewJSONLogger(kitLog.NewSyncWriter(os.Stderr))
	logger = kitLog.With(logger, "ts", kitLog.DefaultTimestampUTC)

	// Service Logger
	sl := kitLog.With(logger, "component", "http")

	// Server instance
	svr := server.NewService(sl, config)

	// -----------   Initializing Service ----------- //

	// User Service
	pkg.InitUserService(db, sl, config, svr)

	// Mio Service
	pkg.InitMioService(db, sl, config, svr)

	// Products Service
	pkg.InitProducts(db, sl, config, svr)

	// -----------   Initializing Service ----------- //

	// listen and serve...
	svr.Listen()
}
