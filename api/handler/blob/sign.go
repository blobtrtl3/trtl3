package handler

import (
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/gin-gonic/gin"
)

type signBlobReq struct {
	Bucket string `json:"bucket" binding:"required,alphanum"`
	ID string `json:"id" binding:"required,alphanum"`
}

// @Summary      Serve a blob
func (bh *BlobHandler) Sign(c *gin.Context) {
	var req signBlobReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
      "message": "verify the data you sent",
    })
		return
	}

	bh.hashmap["abc"] = domain.Signature{
		Bucket: req.Bucket,
		ID: req.ID,
		TTL: time.Now().Add(50 * time.Minute),
		Once: false,
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"url": "https://localhost:7713/b?sign=abc",
	})
}
