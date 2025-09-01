package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (bh *BlobHandler) Delete(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

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

