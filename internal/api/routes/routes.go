package routes

import (
	"time"

	handler "github.com/blobtrtl3/trtl3/internal/api/handler/blob"
	"github.com/blobtrtl3/trtl3/internal/api/middleware"
	"github.com/blobtrtl3/trtl3/internal/cache"
	"github.com/blobtrtl3/trtl3/internal/engine/blob"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/gin-gonic/gin"
)

type RoutesCtx struct {
	r          *gin.Engine
	blobEngine    blob.BlobEngine
	signaturesCache cache.SignaturesCache
	blobQueue  queue.BlobQueue
}

func NewRoutesCtx(r *gin.Engine, be blob.BlobEngine, sc cache.SignaturesCache, q queue.BlobQueue) *RoutesCtx {
	return &RoutesCtx{
		r:          r,
		blobEngine:    be,
		signaturesCache: sc,
		blobQueue:  q,
	}
}

func (rctx *RoutesCtx) SetupRoutes() {
	blobHandler := handler.NewBlob(rctx.blobEngine, rctx.signaturesCache, rctx.blobQueue)

	// Health check endpoint (no authentication required)
	rctx.r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "trtl3",
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
