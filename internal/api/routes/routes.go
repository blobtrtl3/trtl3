package routes

import (
	handler "github.com/blobtrtl3/trtl3/internal/api/handler/blob"
	"github.com/blobtrtl3/trtl3/internal/api/middleware"
	"github.com/blobtrtl3/trtl3/internal/usecase/signatures"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/gin-gonic/gin"
)

type RoutesCtx struct {
	r *gin.Engine
	storage storage.Storage
	signatures signatures.Signatures
}

func NewRoutesCtx(r *gin.Engine, st storage.Storage, si signatures.Signatures) *RoutesCtx {
	return &RoutesCtx{
		r: r,
		storage: st,
		signatures: si,
	}
}

func (rctx *RoutesCtx) SetupRoutes() {
	blobHandler := handler.NewBlob(rctx.storage, rctx.signatures)

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

