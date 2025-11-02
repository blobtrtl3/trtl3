package router

import (
	"context"
	"time"

	"github.com/blobtrtl3/trtl3/internal/blob"
	"github.com/blobtrtl3/trtl3/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RouterCtx struct {
	r               *gin.Engine
	blobService     blob.Service
	redis *redis.Client
	ctx context.Context
}

func NewRouterCtx(r *gin.Engine, bs blob.Service, re *redis.Client) *RouterCtx {
	return &RouterCtx{
		r:               r,
		blobService:     bs,
		redis: re,
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

	serve := rctx.r.Group("/b", middleware.SignMiddleware(rctx.ctx, rctx.redis))
	{
		serve.GET("", handler.Serve)
	}
}
