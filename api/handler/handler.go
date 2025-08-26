package handler

import (
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/gin-gonic/gin"
)

type Blob struct {
	storage storage.Storage
}

func NewBlob(s storage.Storage) *Blob {
	return &Blob{
		storage: s,
	}
}

func (b *Blob) Save(c *gin.Context) {
	c.JSON(200, gin.H{"message": "blob created"})
}
