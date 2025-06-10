package mio

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	vendora "github.com/ppeymann/vendors.git"
	"github.com/ppeymann/vendors.git/config"
	"github.com/ppeymann/vendors.git/models"
	"github.com/ppeymann/vendors.git/server"
)

type handler struct {
	next models.MioService
	conf *config.Configuration
}

func NewHandler(svc models.MioService, conf *config.Configuration, srv *server.Server) models.MioHandler {
	handler := &handler{
		next: svc,
		conf: conf,
	}

	group := srv.Router.Group("/api/v1/storage")
	{
		group.GET("/download/:token", handler.Download)
		group.GET("/image/:size/:token", handler.Image)

		group.Use(srv.Authenticate(vendora.AllRoles))
		{
			group.POST("/upload/:tag", handler.Upload)
		}
	}

	return handler
}

// Download handles downloading files request.
//
// @BasePath 		/api/v1/mio/download
// @Summary			uploading file to mio service
// @Description 	upload specified file to mio service with specified properties
// @Tags 			mio
// @Accept 			octet-stream
// @Produce 		octet-stream
// @Param			token		path		string			true 	"access token of file"
// @Success 		200 		{object} 	vendora.BaseResult	"always returns status 200 but body contains errors"
// @Router 			/mio/download/{token}	[get]
// @Security		Authenticate Header
func (h *handler) Download(ctx *gin.Context) {
	in := &models.DownloadInput{
		Token: ctx.Param("token"),
	}

	data, file, err := h.next.Download(in, ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, vendora.BaseResult{
			Errors: []string{err.Error()},
		})

		return
	}

	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	ctx.Writer.Header().Set("Content-Type", file.ContentType)
	ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	ctx.Writer.Header().Set("Vary", "Accept")

	_, _ = ctx.Writer.Write(data)
}

// Image handles get sized image request.
//
// @BasePath 		/api/v1/mio/image
// @Summary			uploading file to mio service
// @Description 	upload specified file to mio service with specified properties
// @Tags 			mio
// @Accept 			octet-stream
// @Produce 		octet-stream
// @Param			size		path		int		true	"width of requested image"
// @Param			token		path		string	true	"access token of file"
// @Success 		200 		{object} 	vendora.BaseResult	"always returns status 200 but body contains errors"
// @Router 			/mio/image/{size}/{token}	[get]
// @Security		Authenticate Header
func (h *handler) Image(ctx *gin.Context) {
	s, ok := ctx.Params.Get("size")
	if !ok {
		ctx.JSON(http.StatusBadRequest, vendora.BaseResult{
			Errors: []string{"image width is not provided"},
		})
	}

	size, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, vendora.BaseResult{
			Errors: []string{"incompatible image width type"},
		})
	}

	in := &models.DownloadInput{
		Token: ctx.Param("token"),
		Size:  uint(size),
	}

	data, file, err := h.next.Image(in, ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, vendora.BaseResult{
			Errors: []string{err.Error()},
		})

		return
	}

	// call associated service method for expected request
	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	ctx.Writer.Header().Set("Content-Type", file.ContentType)
	ctx.Writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
	ctx.Writer.Header().Set("Vary", "Accept")

	_, _ = ctx.Writer.Write(data)
}

// Upload handler uploading files request.
//
// @BasePath			/api/v1/mio/upload
// @Summary				uploading file to mio service
// @Description			uploading specified file to mio service with specified properties
// @Tags 				mio
// @Accept 				mpfd
// @Produce 			json
// @Param				file			formData	file			true	"uploading file"
// @Param				tag				path		string			true 	"string enums" Enums(public, private, chat, profile)		"uploading file tag"
//
// @Param				Authenticate header string false "authentication paseto token [Required If AuthMode: paseto]"
//
// @Success 			200 		{object} 	vendora.BaseResult		"always returns status 200 but body contains errors"
// @Router 				/mio/upload/{tag}	[post]
// @Security			Authenticate Header
// @Security			Session
func (h *handler) Upload(ctx *gin.Context) {
	in := &models.UploadInput{
		Tag:  ctx.Param("tag"),
		Size: 0,
	}

	result := h.next.Upload(in, ctx)
	ctx.JSON(result.Status, result)
}
