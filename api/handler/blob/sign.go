package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/shared"
	"github.com/gin-gonic/gin"
)

type signBlobReq struct {
	Bucket string `json:"bucket" binding:"required,alphanum"`
	ID     string `json:"id" binding:"required,alphanum"`
	TTL    int    `json:"ttl" binding:"required,min=1,max=24"`
	Once   bool   `json:"once"`
}

// @Summary      Sign a blob url
// @Description  Sign a blob to others acces it without server key
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param        request body signBlobReq true "Blob sign request"
// @Success      201
// @Router       /blobs/sign [post]
func (bh *BlobHandler) Sign(c *gin.Context) {
	var req signBlobReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "verify the data you sent",
		})
		return
	}

	if _, err := bh.storage.FindUnique(req.Bucket, req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "verify the data you sent and try again",
		})
		return
	}

	now := time.Now()
	signature := fmt.Sprintf("%s%s", shared.GenShortID(), now.Format("050204")) // format to SSDDMM

	bh.signatures.Set(
		signature, 
		domain.Signature{
			Bucket: req.Bucket,
			ID:     req.ID,
			TTL:    now.Add(time.Duration(req.TTL) * time.Hour),
			Once:   req.Once,
		},
	)

	c.JSON(http.StatusCreated, gin.H{
		"url": fmt.Sprintf("https://localhost:7713/b?sign=%s", signature),
	})
}
