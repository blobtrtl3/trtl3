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

type saveBlobReq struct {
	Bucket string `form:"bucket" binding:"required,alphanum"`
}

type BlobHandler struct {
	storage storage.Storage
}

func NewBlob(s storage.Storage) *BlobHandler {
	return &BlobHandler{
		storage: s,
	}
}

func (bh *BlobHandler) Save(c *gin.Context) {
	var req saveBlobReq

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the data you sent"})
		return
	}

	blobMultipart, err := c.FormFile("blob")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not open the blob that you sent"})
		return
	}

	inMemBlob, err := blobMultipart.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorry, we had an error, try again"})
		return
	}
	defer inMemBlob.Close()

	blobBytes, err := io.ReadAll(inMemBlob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorry, we had an error, try again"})
		return
	}

	blobInfo := &domain.BlobInfo{
		ID:        shared.GenShortID(),
		Bucket:    req.Bucket,
		Mime:      blobMultipart.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
		Size:      int(int64(blobMultipart.Size) / 1024), // NOTE: blob.Size return value in bytes so I did it to be an KB value
	}

	_, err = bh.storage.Save(blobInfo, &blobBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the blob, try again"})
		return
	}

	c.JSON(200, blobInfo)
}

func (bh *BlobHandler) FindByBucketOrID(c *gin.Context) {
	bucket := c.Query("bucket")
	id := c.Query("id")

	if id == "" && bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	if id == "" && bucket != "" { // then find only by bucket
		blobsInfos, err := bh.storage.FindByBucket(bucket)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s", bucket)})
			return
		}

		if blobsInfos == nil {
			c.JSON(http.StatusOK, gin.H{"blobs": ""})
			return
		}

		c.JSON(http.StatusOK, gin.H{"blobs": blobsInfos})
		return
	}
	// here find by bucket and id

	blobInfo, err := bh.storage.FindByBucketAndID(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, blobInfo)
}

func (bh *BlobHandler) Delete(c *gin.Context) {
	bucket := c.Query("bucket")
	id := c.Query("id")

	if bucket == "" && id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	_, err := bh.storage.Delete(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob by in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blob deleted"})
}

func (bh *BlobHandler) DownloadByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the id sent"})
		return
	}

	blobBytes, err := bh.storage.DownloadByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob with id: %s", id)})
		return
	}

	c.JSON(http.StatusOK, blobBytes)
}
