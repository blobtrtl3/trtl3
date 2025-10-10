package router

import (
	"time"

	"github.com/blobtrtl3/trtl3/internal/http/handler"
	"github.com/blobtrtl3/trtl3/internal/http/middleware"
	"github.com/blobtrtl3/trtl3/internal/infra/cache"
	"github.com/blobtrtl3/trtl3/internal/service"
	"github.com/gin-gonic/gin"
)

type RouterCtx struct {
	r               *gin.Engine
	BlobService     service.BlobService
	signaturesCache cache.SignaturesCache
}

func NewRouterCtx(r *gin.Engine, bs service.BlobService, sc cache.SignaturesCache) *RouterCtx {
	return &RouterCtx{
		r:               r,
		BlobService:     bs,
		signaturesCache: sc,
	}
}

func (rctx *RouterCtx) SetupRouter() {
	blobHandler := handler.NewBlobHandler(rctx.BlobService)

	// Health check endpoint (no authentication required)
	rctx.r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "trtl3",
			"timestamp": time.Now().UTC(),
		})
	})

	protected := rctx.r.Group("/blobs", middleware.AuthMiddleware())
	{
		protected.POST("", blobHandler.Save)
		protected.GET("", blobHandler.FindByBucket)
		protected.GET("/:bucket/:id", blobHandler.FindUnique)
		protected.DELETE("/:bucket/:id", blobHandler.Delete)

		protected.GET("/download/:bucket/:id", blobHandler.Download)

		protected.POST("/sign", blobHandler.Sign)
	}

	serve := rctx.r.Group("/b", middleware.SignMiddleware(rctx.signaturesCache))
	{
		serve.GET("", blobHandler.Serve)
	}
}
