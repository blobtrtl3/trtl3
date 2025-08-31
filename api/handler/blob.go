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
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the data you sent"})
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

func (b *Blob) FindByBucketOrID(c *gin.Context) {
	id := c.Query("id")
  bucket := c.Query("bucket")

  if id != "" { // then it will find using id
		blob, err := b.storage.FindByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob by with id: %s", id)})
			return
		}

		c.JSON(http.StatusOK, blob)
		return
  }

	if bucket != "" { // then it will find using bucket
		blobs, err := b.storage.FindByBucket(bucket)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s", bucket)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"blobs": blobs})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id"})
}

func (b *Blob) Delete(c *gin.Context) {
	id := c.Query("id")

  if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the id sent"})
		return
  }

	if err := b.storage.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob by with id: %s", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blob deleted"})
}
