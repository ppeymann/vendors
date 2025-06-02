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
// @Success 					200 {object} models.BaseResult{result=models.TokenBundlerOutput}
// @Router						POST /api/v1/user/register
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
