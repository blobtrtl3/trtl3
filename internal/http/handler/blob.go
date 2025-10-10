package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blobtrtl3/trtl3/internal/service"
	"github.com/gin-gonic/gin"
)

type BlobHandler struct {
	blobService service.BlobService
}

func NewBlobHandler(bs service.BlobService) *BlobHandler {
	return &BlobHandler{blobService: bs}
}

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

	fileBlob, err := blobMultipart.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorry, we had an error, try again"})
		return
	}
	defer fileBlob.Close()

	blobInfo, err := bh.blobService.Save(
		req.Bucket,
		blobMultipart.Header.Get("Content-Type"),
		blobMultipart.Size, // NOTE: size in bytes value
		fileBlob,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "sorry, we had an error, try again"})
		return
	}

	c.JSON(http.StatusCreated, blobInfo)
}

// @Summary      Find blobs by bucket
// @Description  Find all blobs inside a bucket
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param        bucket query string true "Bucket name"
// @Success      200 {object} []domain.BlobInfo
// @Router       /blobs [get]
func (bh *BlobHandler) FindByBucket(c *gin.Context) {
	bucket := c.Query("bucket")

	blobsInfos, err := bh.blobService.FindByBucket(bucket)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s", bucket)})
		return
	}

	if blobsInfos == nil {
		c.JSON(http.StatusOK, gin.H{"blobs": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blobs": blobsInfos})
}

// @Summary      Find unique blob
// @Description  Find blob by using bucket and id
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param 			 bucket path string true "Bucket name"
// @Param 			 id path string true "Blob id"
// @Success      200 {object} domain.BlobInfo
// @Router       /blobs/{bucket}/{id} [get]
func (bh *BlobHandler) FindUnique(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

	blobInfo, err := bh.blobService.FindUnique(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, blobInfo)
}

// @Summary      Download blob
// @Description  Download a single blob by using the bucket and id
// @Accept       json
// @Produce      application/octet-stream
// @Param        Authorization header string true "Access token" default("")
// @Param 			 bucket path string true "Bucket name"
// @Param 			 id path string true "Blob id"
// @Success 		 200 {file} file "Binary file"
// @Router       /blobs/download/{bucket}/{id} [get]
func (bh *BlobHandler) Download(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

	blobBytes, err := bh.blobService.Download(bucket, id)
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

// @Summary      Delete blob
// @Description  Delete a single blob by using the bucket and id
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Access token" default("")
// @Param 			 bucket path string true "Bucket name"
// @Param 			 id path string true "Blob id"
// @Success      200
// @Router       /blobs/{bucket}/{id} [get]
func (bh *BlobHandler) Delete(c *gin.Context) {
	bucket := c.Param("bucket")
	id := c.Param("id")

	_, err := bh.blobService.Delete(bucket, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("could not find blob by in bucket: %s with id: %s", bucket, id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "blob deleted"})
}

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

	serveInfo, err := bh.blobService.Serve(bucket.(string), id.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "we had an error, sorry try again"})
		return
	}

	c.Header("Content-Type", serveInfo.Mime)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", time.Now().Format("020504"))) //DDSSMM

	c.File(serveInfo.Path)
}

type signBlobReq struct {
	Bucket string `json:"bucket" binding:"required,alphanum"`
	ID     string `json:"id" binding:"required,alphanum"`
	TTL    int    `json:"ttl" binding:"required,min=1,max=1440"`
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

	signature, err := bh.blobService.Sign(
		req.Bucket,
		req.ID,
		req.TTL,
		req.Once,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "verify the data you sent and try again",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"url": fmt.Sprintf("http://localhost:7713/b?sign=%s", signature),
	})
}
