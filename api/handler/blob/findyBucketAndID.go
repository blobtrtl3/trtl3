package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Find unique blob
// @Description  Find blob by using bucket and id
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param 			 bucket path string true "Bucket name"
// @Param 			 id path string true "Blob id"
// @Success      200 {object} domain.BlobInfo
// @Router       /blobs/{bucket}/{id} [get]
func (bh *BlobHandler) FindByBucketAndID(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

	if id == "" && bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	blobInfo, err := bh.storage.FindByBucketAndID(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, blobInfo)
}
