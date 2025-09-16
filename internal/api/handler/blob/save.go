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

// @Summary      Upload blob
// @Description  Upload a blob to server
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param        bucket formData string true "Bucket name"
// @Param        blob formData file true "Blob file"
// @Success      201 {object} domain.BlobInfo
// @Router       /blobs [post]
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
		Size:      blobMultipart.Size, // NOTE: size in bytes value
	}

	_, err = bh.storage.Save(blobInfo, blobBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the blob, try again"})
		return
	}

	c.JSON(http.StatusCreated, blobInfo)
}
