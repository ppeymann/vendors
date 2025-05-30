package server

import (
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/log"
	"github.com/ppeymann/vendors.git/auth"
	"github.com/ppeymann/vendors.git/config"
)

type Server struct {
	Router        *gin.Engine
	Config        *config.Configuration
	Logger        kitlog.Logger
	paseto        auth.TokenMaker
	instrumenting serviceInstrumenting
}
