package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/shared"
	"github.com/gin-gonic/gin"
)

type saveBlobReq struct {
	Bucket string `form:"bucket" binding:"required,alphanum"`
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

	_, err = bh.storage.Save(blobInfo, blobBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the blob, try again"})
		return
	}

	c.JSON(200, blobInfo)
}
