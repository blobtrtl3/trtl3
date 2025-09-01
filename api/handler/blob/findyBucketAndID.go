package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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


