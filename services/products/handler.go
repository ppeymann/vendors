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

// DeleteProduct handler request
//
// @BasePath			/api/v1/product
// @Summary				delete a product
// @Description			delete a product with specific ID
// @Tags				products
// @Accept				json
// @Produce				json
//
// @Param				id		path	string		true		"product ID"
// @Success				200		{object}	vendora.BaseResult{result=uint}		"always return status 200 but body contains error"
// @Router				/{id}	[delete]
// @Security			session
func (h *handler) DeleteProduct(ctx *gin.Context) {
	id, err := server.GetPathUint64(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredParam},
		})

		return
	}

	result := h.next.DeleteProduct(ctx, uint(id))
	ctx.JSON(result.Status, result)
}

// EditProduct handler http request.
//
// @BasePath			/api/v1/products
// @Summary				edit a product
// @Description			edit a product with specific ID
// @Tags				products
// @Accept				json
// @Produce				json
//
// @Param				input	body	models.ProductInput		true	"product input"
// @Success				200		{object}	vendora.BaseResult{result=models.ProductEntity}		"always return status 200 but body contains error"
// @Router				/edit/{id}	[patch]
// @Security			session
func (h *handler) EditProduct(ctx *gin.Context) {
	in := &models.ProductInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredJsonBody},
		})

		return
	}

	id, err := server.GetPathUint64(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredParam},
		})

		return
	}

	result := h.next.EditProduct(ctx, uint(id), in)
	ctx.JSON(result.Status, result)
}

// GetByTags handler http request
//
// @BasePath			/api/v1/products
// @Summary				get products with same tags
// @Tags				products
// @Accept				json
// @Produce				json
//
// @Param				input	 body			models.TagsInput		true		"slice of tags"
// @Success				200		{object}		vendora.BaseResult{result=[]models.ProductEntity}		"always return status 200 but body contains error"
// @Router				/tags	[post]
func (h *handler) GetByTags(ctx *gin.Context) {
	in := &models.TagsInput{}
	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.GetByTags(ctx, in.Tags)
	ctx.JSON(result.Status, result)
}

// GetProduct is handler for get product
//
// @BasePath			/api/v1/products
// @Summary				GetProduct with specific ID
// @Tags				products
// @Accept				json
// @Product				json
//
// @Param				id		path		string	true	"product ID"
// @Success				200		{object}	vendora.BaseResult{result=models.ProductEntity}		"always returns status 200 but body contains error"
// @Router				/{id}	[get]
func (h *handler) GetProduct(ctx *gin.Context) {
	id, err := server.GetPathUint64(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &vendora.BaseResult{
			Errors: []string{vendora.ProvideRequiredParam},
		})

		return
	}

	result := h.next.GetProduct(ctx, uint(id))
	ctx.JSON(http.StatusOK, result)

}

// Add New Product
//
// @BasePath			/api/v1/products
// @Summary				Add New Product
// @Description			Add new production with specific User
// @Tags				products
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

func NewHandler(srv models.ProductService, s *server.Server) models.ProductHandler {
	handler := &handler{
		next: srv,
	}

	group := s.Router.Group("/api/v1/product")
	{
		group.GET("/:id", handler.GetProduct)
		group.GET("/tags", handler.GetByTags)
	}
	group.Use(s.Authenticate(vendora.SellerRoles))
	{
		group.POST("/add", handler.Add)
		group.PATCH("/edit/:id", handler.EditProduct)
	}

	return handler
}
