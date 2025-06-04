package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/auth"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/models"

	"github.com/ppeymann/vendors.git/utils"
)

type service struct {
	repo models.UserRepository
	conf *config.Configuration
}

func NewService(repo models.UserRepository, conf *config.Configuration) models.UserService {
	return &service{
		repo: repo,
		conf: conf,
	}
}

func (s *service) Register(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	if env.IsProduction() {
		pass, err := utils.HashString(in.Password)
		if err != nil {
			return &vendora.BaseResult{
				Errors: []string{err.Error()},
				Status: http.StatusOK,
			}
		}

		in.Password = pass
	}

	user, err := s.repo.Create(in)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	if s.conf.Listener.AuthMode == config.Session {
		session := sessions.Default(ctx)
		session.Set(vendora.UserSessionKey, auth.Claims{
			Subject: uint(user.ID),
			Roles:   user.Roles,
		})

		err = session.Save()
		if err != nil {
			return &vendora.BaseResult{
				Status: http.StatusOK,
				Errors: []string{vendora.ErrInternalServer.Error()},
			}
		}

		return &vendora.BaseResult{
			Status: http.StatusOK,
			Result: "you are Logged In Success",
		}
	}
	return nil
}
