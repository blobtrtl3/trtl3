package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/blobtrtl3/trtl3/shared"
	"github.com/gin-gonic/gin"
)

type saveBlobRequest struct {
	Bucket string `form:"bucket" binding:"required,alphanum"`
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

	if err := c.ShouldBind(&br); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the data you sent"})
		return
	}

	blob, err := c.FormFile("blob")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not open the blob that you sent"})
		return
	}

	inMemBlob, err := blob.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorry, we had an error, try again"})
		return
	}
	defer inMemBlob.Close()

	bodyBytes, err := io.ReadAll(inMemBlob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorry, we had an error, try again"})
		return
	}

	bi := &domain.BlobInfo{
		ID:        shared.GenShortID(),
		Bucket:    br.Bucket,
		Mime:      blob.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
		Size:      int(int64(blob.Size) / 1024), // NOTE: blob.Size return value in bytes so I did it to be an KB value
	}

	_, err = b.storage.Save(bi, &bodyBytes)
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
