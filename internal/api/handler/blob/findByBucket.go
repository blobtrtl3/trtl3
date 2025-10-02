package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Find blobs by bucket
// @Description  Find all blobs inside a bucket
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param        bucket query string true "Bucket name"
// @Success      200 {object} []domain.BlobInfo
// @Router       /blobs [get]
func (bh *BlobHandler) FindByBucket(c *gin.Context) {
	bucket := c.Query("bucket")

	if bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket sent"})
		return
	}

	blobsInfos, err := bh.blobEngine.FindByBucket(bucket)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s", bucket)})
		return
	}

	if blobsInfos == nil {
		c.JSON(http.StatusOK, gin.H{"blobs": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blobs": blobsInfos})
}
