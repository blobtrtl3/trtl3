package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	//TODO ta enviando como string
	c.JSON(http.StatusOK, blobBytes)
}
