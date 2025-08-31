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

	_, err := b.storage.Save(bi, &bbytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the blob, try again"})
		return
	}

	c.JSON(200, bi)
}

func (b *Blob) FindByBucketOrID(c *gin.Context) {
  bucket := c.Query("bucket")
	id := c.Query("id")

  if id == "" && bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
  }

	if id == "" && bucket != "" { // then find only by bucket
		blobs, err := b.storage.FindByBucket(bucket)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s", bucket)})
			return
		}

		if blobs == nil {
			c.JSON(http.StatusOK, gin.H{"blobs": ""})
			return
		}

		c.JSON(http.StatusOK, gin.H{"blobs": blobs})
		return
	}
	// here find by bucket and id

	blob, err := b.storage.FindByBucketAndID(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, blob)
}

func (b *Blob) Delete(c *gin.Context) {
	bucket := c.Query("bucket")
	id := c.Query("id")

  if bucket == "" && id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
  }

	_, err := b.storage.Delete(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob by in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blob deleted"})
}
