package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary      Download blob
// @Description  Download a single blob by using the bucket and id
// @Accept       json
// @Produce      application/octet-stream
// @Param        Authorization header string true "Access token" default("")
// @Param 			 bucket path string true "Bucket name"
// @Param 			 id path string true "Blob id"
// @Success      200 {object}
// @Router       /blobs/download/{bucket}/{id} [get]
func (bh *BlobHandler) DownloadByID(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

	if bucket == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	blobBytes, err := bh.storage.DownloadByID(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.Data(
		http.StatusOK,
		"application/octet-stream",
		blobBytes,
	)

}
