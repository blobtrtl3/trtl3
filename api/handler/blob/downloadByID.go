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
	id := c.Param("id")
	// TODO: take by id and bucket

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the id sent"})
		return
	}

	blobBytes, err := bh.storage.DownloadByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob with id: %s", id)})
		return
	}

	c.Data(
		http.StatusOK,
		"application/octet-stream",
		blobBytes,
	)

}
