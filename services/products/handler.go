package products

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/server"
)

type handler struct {
	next models.ProductService
}

func NewHandler(srv models.ProductService, s *server.Server) models.ProductHandler {
	handler := &handler{
		next: srv,
	}

	group := s.Router.Group("/api/v1/product")
	{

	}
	group.Use(s.Authenticate(vendora.SellerRoles))
	{
		group.POST("/add", handler.Add)
	}

	return handler
}

// Add New Product
//
// @BasePath			/api/v1/product
// @Summary				Add New Product
// @Description			Add new production with specific User
// @Accept				json
// @Product				json
//
// @Param				input	body	models.ProductInput	true	"Product Input"
// @Success				200		{object}	vendora.BaseResult{result=models.ProductEntity}		"always returns status 200 but body contains error"
// @Router				/add	[post]
// @Security			session
func (h *handler) Add(ctx *gin.Context) {
	in := &models.ProductInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.Add(ctx, in)
	ctx.JSON(result.Status, result)
}
