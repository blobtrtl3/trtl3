package routes

import (
	"time"

	handler "github.com/blobtrtl3/trtl3/internal/api/handler/blob"
	"github.com/blobtrtl3/trtl3/internal/api/middleware"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/blobtrtl3/trtl3/internal/repo/signatures"
	"github.com/blobtrtl3/trtl3/internal/repo/storage"
	"github.com/gin-gonic/gin"
)

type RoutesCtx struct {
	r          *gin.Engine
	storage    storage.Storage
	signatures signatures.Signatures
	blobQueue  queue.BlobQueue
}

func NewRoutesCtx(r *gin.Engine, st storage.Storage, si signatures.Signatures, q queue.BlobQueue) *RoutesCtx {
	return &RoutesCtx{
		r:          r,
		storage:    st,
		signatures: si,
		blobQueue:  q,
	}
}

func (rctx *RoutesCtx) SetupRoutes() {
	blobHandler := handler.NewBlob(rctx.storage, rctx.signatures, rctx.blobQueue)

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

	serve := rctx.r.Group("/b", middleware.SignMiddleware(rctx.signatures))
	{
		serve.GET("", blobHandler.Serve)
	}
}
