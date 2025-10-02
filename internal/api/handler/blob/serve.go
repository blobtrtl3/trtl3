package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/shared"
	"github.com/gin-gonic/gin"
)

// @Summary      Serve a blob
// @Description  Serve a blob by a signed url
// @Accept       json
// @Produce      json
// @Param 			 sign path string true "sign"
// @Success      200
// @Router       /b [get]
func (bh *BlobHandler) Serve(c *gin.Context) {
	bucket, exists := c.Get("bucket")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "blob not found"})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "blob not found"})
		return
	}

	blobInfo, err := bh.blobEngine.FindUnique(bucket.(string), id.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "blob not found, check your data and try again"})
		return
	}

	blobName := shared.GenBlobName(blobInfo.Bucket, blobInfo.ID)

	c.Header("Content-Type", blobInfo.Mime)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", time.Now().Format("020504"))) //DDSSMM

	path := filepath.Join("blobs", blobName)
	c.File(path)
}
