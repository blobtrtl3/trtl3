package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Delete blob
// @Description  Delete a single blob by using the bucket and id
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param 			 bucket path string true "Bucket name"
// @Param 			 id path string true "Blob id"
// @Success      200
// @Router       /blobs/{bucket}/{id} [get]
func (bh *BlobHandler) Delete(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

	if bucket == "" && id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	_, err := bh.blobEngine.Delete(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob by in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blob deleted"})
}
