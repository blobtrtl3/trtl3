package router

import (
	"time"

	"github.com/blobtrtl3/trtl3/internal/blob"
	"github.com/blobtrtl3/trtl3/internal/http/middleware"
	"github.com/blobtrtl3/trtl3/internal/infra"
	"github.com/gin-gonic/gin"
)

type RouterCtx struct {
	r               *gin.Engine
	blobService     blob.Service
	signaturesCache infra.SignaturesCache
}

func NewRouterCtx(r *gin.Engine, bs blob.Service, sc infra.SignaturesCache) *RouterCtx {
	return &RouterCtx{
		r:               r,
		blobService:     bs,
		signaturesCache: sc,
	}
}

func (rctx *RouterCtx) SetupRouter() {
	handler := blob.NewHandler(rctx.blobService)

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
		protected.POST("", handler.Save)
		protected.GET("", handler.FindByBucket)
		protected.GET("/:bucket/:id", handler.FindUnique)
		protected.DELETE("/:bucket/:id", handler.Delete)

		protected.GET("/download/:bucket/:id", handler.Download)

		protected.POST("/sign", handler.Sign)
	}

	serve := rctx.r.Group("/b", middleware.SignMiddleware(rctx.signaturesCache))
	{
		serve.GET("", handler.Serve)
	}
}
