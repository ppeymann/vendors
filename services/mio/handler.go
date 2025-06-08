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

// Download implements models.MioHandler.
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

// Image implements models.MioHandler.
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

// Upload implements models.MioHandler.
func (h *handler) Upload(ctx *gin.Context) {
	in := &models.UploadInput{
		Tag:  ctx.Param("tag"),
		Size: 0,
	}

	result := h.next.Upload(in, ctx)
	ctx.JSON(result.Status, result)
}
