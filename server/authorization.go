package server

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/auth"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/utils"
	"github.com/thoas/go-funk"
)

func (s *Server) initSession() error {
	db, err := sql.Open("postgres", env.GetEnv("DSN", ""))
	if err != nil {
		return err
	}

	store, err := postgres.NewStore(db, []byte(s.Config.Listener.SessionsSecret))
	if err != nil {
		return err
	}

	gob.Register(auth.Claims{})
	s.Router.Use(sessions.Sessions(vendora.UserSessionKey, store))
	return nil
}

// sessionAuth is authentication and authorization middleware for http request
// by passing roles string slice it will control whether context has expected roles or no.
// if expected roles is not important, the roles arg can be empty slice.
func (s *Server) sessionAuth(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		su := session.Get(vendora.UserSessionKey)
		if su == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		user := auth.Claims{}
		err := mapstructure.Decode(su, &user)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		// control access roles if it expected roles not empty
		if len(roles) > 0 {
			if len(funk.IntersectString(roles, user.Roles)) == 0 {
				ctx.AbortWithStatus(http.StatusUnauthorized)

				return
			}
		}

		ctx.Set(vendora.ContextUserKey, user)
		ctx.Next()
	}
}

// Authenticate is authentication and Authenticate middleware for http request
func (s *Server) pasetoAuth(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// catch Authenticate header from context
		ah := ctx.GetHeader("Authorization")

		// abort request if Authenticate header is empty or not provided.
		if len(ah) == 0 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("authorization header is not provided"))
			return
		}

		// Bearer token format validation
		fields := strings.Fields(ah)
		if len(fields) != 2 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("invalid Authorization header format"))
			return
		}

		at := strings.ToLower(fields[0])
		if at != "bearer" {
			_ = ctx.AbortWithError(http.StatusUnauthorized,
				fmt.Errorf("unsupported Authenticate format : %s", fields[0]))
			return
		}

		token := fields[1]
		claims, err := s.paseto.VerifyToken(token)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if claims.ExpiredAt.Before(time.Now().UTC()) {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("authorization token is expired"))
			return
		}

		if len(roles) > 0 {
			if len(funk.IntersectString(roles, claims.Roles)) == 0 {
				_ = ctx.AbortWithError(http.StatusForbidden, errors.New("permission denied"))
				return
			}
		}

		ctx.Set(utils.ContextUserKey, claims)
		ctx.Set(utils.ContextRoleKey, claims.Roles)
		ctx.Next()
	}
}

// Authenticate is authentication and Authenticate middleware for http request
func (s *Server) Authenticate(roles []string) gin.HandlerFunc {
	if s.Config.Listener.AuthMode == config.Session {
		return s.sessionAuth(roles)
	}

	return s.pasetoAuth(roles)
}
