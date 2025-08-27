package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/blobtrtl3/trtl3/shared"
	"github.com/gin-gonic/gin"
)

type saveBlobRequest struct {
	Bucket string `json:"bucket" binding:"required,alphanum"`
}

type Blob struct {
	storage storage.Storage
}

func NewBlob(s storage.Storage) *Blob {
	return &Blob{
		storage: s,
	}
}

func (b *Blob) Save(c *gin.Context) {
	var br saveBlobRequest
	var bbytes []byte // TODO: take blob of request

	if err := c.ShouldBindJSON(&br); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "verirify the data you sent"})
		return
	}

	bi := &domain.BlobInfo{
		ID:        shared.GenShortID(),
		Bucket:    br.Bucket,
		Mime:      "text/plain", // TODO: take mime from headers
		CreatedAt: time.Now(),
		Size:      24, // TODO: calc blob size
	}

	b.storage.Save(bi, &bbytes)

	c.JSON(200, gin.H{"message": "blob created"})
}

func (b *Blob) FindByID(c *gin.Context) {
	c.JSON(200, gin.H{"message": "blob created"})
}

func (b *Blob) FindByBucket(c *gin.Context) {
	c.JSON(200, gin.H{"message": "blob created"})
}

func (b *Blob) Delete(c *gin.Context) {
	c.JSON(200, gin.H{"message": "blob created"})
}
