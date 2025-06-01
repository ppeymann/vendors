package user

import (
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/server"
)

type handler struct {
	next models.UserService
}

func NewHandler(srv models.UserService, s *server.Server) models.UserHandler {
	handler := &handler{
		next: srv,
	}

	group := s.Router.Group("/api/v1/user")
	{
		group.POST("/signup")
	}

	return handler
}
