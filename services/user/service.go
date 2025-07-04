package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/auth"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/env"
	"github.com/ppeymann/vendors.git/models"
	"github.com/segmentio/ksuid"

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

		err = s.repo.Update(user)
		if err != nil {
			return &vendora.BaseResult{
				Errors: []string{err.Error()},
				Status: http.StatusOK,
			}
		}

		return &vendora.BaseResult{
			Status: http.StatusOK,
			Result: "you are Logged In Success",
		}
	}

	referesh := models.RefreshTokenEntity{
		TokenId:   ksuid.New().String(),
		UserAgent: ctx.Request.UserAgent(),
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiredAt: time.Now().Add(time.Duration(s.conf.Jwt.RefreshExpire) * time.Minute).UTC().Unix(),
	}

	user.Tokens = append(user.Tokens, referesh)

	err = s.repo.Update(user)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	// create token and refresh token
	paseto, err := auth.NewPasetoMaker(env.GetEnv("JWT", ""))
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	tokenClaims := &auth.Claims{
		Subject:   user.ID,
		Issuer:    s.conf.Jwt.Issuer,
		Audience:  s.conf.Jwt.Audience,
		IssuedAt:  time.Unix(referesh.IssuedAt, 0),
		ExpiredAt: time.Now().Add(time.Duration(s.conf.Jwt.TokenExpire) * time.Minute).UTC(),
	}

	refereshClaims := &auth.Claims{
		Subject:   user.ID,
		ID:        referesh.TokenId,
		Issuer:    s.conf.Jwt.Issuer,
		Audience:  s.conf.Jwt.Audience,
		IssuedAt:  time.Unix(referesh.IssuedAt, 0),
		ExpiredAt: time.Unix(referesh.ExpiredAt, 0),
	}

	tokenStr, err := paseto.CreateToken(tokenClaims)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	refereshStr, err := paseto.CreateToken(refereshClaims)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: models.TokenBundlerOutput{
			Token:   tokenStr,
			Refresh: refereshStr,
			Expire:  tokenClaims.ExpiredAt,
		},
	}
}

// Login implements models.UserService.
func (s *service) Login(ctx *gin.Context, in *models.AuthInput) *vendora.BaseResult {
	user, err := s.repo.FindByUserName(in.UserName)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	if env.IsProduction() {
		ok := utils.CheckHashedString(in.Password, user.Password)
		if !ok {
			return &vendora.BaseResult{
				Errors: []string{"username or password is wrong"},
				Status: http.StatusOK,
			}
		}
	} else {
		if in.Password != user.Password {
			return &vendora.BaseResult{
				Errors: []string{"username or password is wrong"},
				Status: http.StatusOK,
			}
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
		err = s.repo.Update(user)
		if err != nil {
			return &vendora.BaseResult{
				Errors: []string{err.Error()},
				Status: http.StatusOK,
			}
		}

		return &vendora.BaseResult{
			Status: http.StatusOK,
			Result: "you are Logged In Success",
		}
	}

	referesh := models.RefreshTokenEntity{
		TokenId:   ksuid.New().String(),
		UserAgent: ctx.Request.UserAgent(),
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiredAt: time.Now().Add(time.Duration(s.conf.Jwt.RefreshExpire) * time.Minute).UTC().Unix(),
	}

	user.Tokens = append(user.Tokens, referesh)

	err = s.repo.Update(user)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	// create token and refresh token
	paseto, err := auth.NewPasetoMaker(env.GetEnv("JWT", ""))
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	tokenClaims := &auth.Claims{
		Subject:   user.ID,
		Issuer:    s.conf.Jwt.Issuer,
		Audience:  s.conf.Jwt.Audience,
		IssuedAt:  time.Unix(referesh.IssuedAt, 0),
		ExpiredAt: time.Now().Add(time.Duration(s.conf.Jwt.TokenExpire) * time.Minute).UTC(),
	}

	refereshClaims := &auth.Claims{
		Subject:   user.ID,
		ID:        referesh.TokenId,
		Issuer:    s.conf.Jwt.Issuer,
		Audience:  s.conf.Jwt.Audience,
		IssuedAt:  time.Unix(referesh.IssuedAt, 0),
		ExpiredAt: time.Unix(referesh.ExpiredAt, 0),
	}

	tokenStr, err := paseto.CreateToken(tokenClaims)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	refereshStr, err := paseto.CreateToken(refereshClaims)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: models.TokenBundlerOutput{
			Token:   tokenStr,
			Refresh: refereshStr,
			Expire:  tokenClaims.ExpiredAt,
		},
	}
}

// User implements models.UserService.
func (s *service) User(ctx *gin.Context) *vendora.BaseResult {
	claims, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	user, err := s.repo.FindByID(claims.Subject)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: user,
	}
}

// EditUser implements models.UserService.
func (s *service) EditUser(ctx *gin.Context, in *models.EditUserInput) *vendora.BaseResult {
	claims, err := vendora.CheckAuth(ctx)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	user, err := s.repo.EditUser(claims.Subject, in)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: user,
	}
}

// GetAllUserWithRole implements models.UserService.
func (s *service) GetAllUserWithRole(_ *gin.Context, role string) *vendora.BaseResult {
	users, err := s.repo.GetAllUserWithRole(role)
	if err != nil {
		return &vendora.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	if len(users) == 0 {
		return &vendora.BaseResult{
			Errors: []string{errors.New("EOF").Error()},
			Status: http.StatusOK,
		}
	}

	return &vendora.BaseResult{
		Status: http.StatusOK,
		Result: users,
	}
}
