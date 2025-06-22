package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
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
		group.POST("/signup", handler.Register)
		group.POST("/login", handler.Login)
	}

	group.Use(s.Authenticate(vendora.AllRoles))
	{
		group.GET("", handler.User)
		group.PATCH("", handler.EditUser)
		group.GET("/:role", handler.GetAllUserWithRole)
	}

	return handler
}

// Register is handler for create New user
//
// @BasePath			 		/api/v1/user
// @Summary		 				Create New user
// @Description		 			Create New user
// @Tags 						user
// @Accept						json
// @Produce 					json
//
// @Param						input body models.AuthInput true "AuthInput"
// @Success 					200 {object} vendora.BaseResult{result=models.TokenBundlerOutput}
// @Router						/api/v1/user/register	[post]
func (h *handler) Register(ctx *gin.Context) {
	in := &models.AuthInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.Register(ctx, in)
	ctx.JSON(result.Status, result)
}

// Login is handler for log in
//
// @BasePath			/api/v1/user
// @Summary				log in
// @Description			log in with specific user name
// @Tags				user
// @Accept				json
// @Produce				json
//
// @Params				input body models.AuthInput	true	"AuthInput"
// @Success				200 {object} vendora.BaseResult{result=models.TokenBundlerOutput}
// @Router				/api/v1/user/login	[post]
func (h *handler) Login(ctx *gin.Context) {
	in := &models.AuthInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.Login(ctx, in)
	ctx.JSON(result.Status, result)
}

// User is handler for get information
//
// @BasePath			/api/v1/user
// @Summary				user info
// @Description			get user information
// @Tags				user
// @Accept				json
// @Produce				json
//
// @Success				200	{object}	vendora.BaseResult{result=models.UserEntity}
// @Router				/api/v1/user	[get]
func (h *handler) User(ctx *gin.Context) {

	result := h.next.User(ctx)
	ctx.JSON(result.Status, result)
}

// EditUser is handler for editing user information like first name , ...
//
// @BasePath			/api/v1/user
// @Summary				edit user
// @Description			edit user with specific ID
// @Tags				user
// @Accept				json
// @Product				json
//
// @Success				200	{object}	vendora.BaseResult{result=models.UserEntity}
// @Router				/api/v1/user/EditUser
func (h *handler) EditUser(ctx *gin.Context) {
	in := &models.EditUserInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.EditUser(ctx, in)
	ctx.JSON(result.Status, result)
}

// GetAllUserWithRole is handler for get all user for admin
//
// @BasePath			/api/v1/user
// @Summary				get all user
// @Description			get all user for admin
// @Tags				user
// @Accept				json
// @Product				json
//
// @Success				200	{object}	vendora.BaseResult{result=[]models.UserEntity}
// @Router				/api/v1/user/{role}	[get]
// @Security			session
func (h *handler) GetAllUserWithRole(ctx *gin.Context) {
	role, err := server.GetStringPath("role", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredParam},
		})

		return
	}

	result := h.next.GetAllUserWithRole(ctx, role)
	ctx.JSON(result.Status, result)
}
