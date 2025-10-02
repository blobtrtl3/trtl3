package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/blobtrtl3/trtl3/internal/cache"
	"github.com/blobtrtl3/trtl3/internal/domain"
	"github.com/blobtrtl3/trtl3/internal/engine"
	"github.com/blobtrtl3/trtl3/internal/queue"
	"github.com/blobtrtl3/trtl3/shared"
	"github.com/gin-gonic/gin"
)

type BlobHandler struct {
	blobEngine      engine.BlobEngine
	signaturesCache cache.SignaturesCache
	bloQueue        queue.BlobQueue
}

func NewBlobHandler(
	be engine.BlobEngine,
	sc cache.SignaturesCache,
	bq queue.BlobQueue,
) *BlobHandler {
	return &BlobHandler{
		blobEngine:      be,
		signaturesCache: sc,
		bloQueue:        bq,
	}
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

	blobInfo := &domain.BlobInfo{
		ID:        shared.GenShortID(),
		Bucket:    req.Bucket,
		Mime:      blobMultipart.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
		Size:      blobMultipart.Size, // NOTE: size in bytes value
	}

	if err = bh.bloQueue.Append(blobInfo, fileBlob); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save the blob, try again"})
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

	if bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket sent"})
		return
	}

	blobsInfos, err := bh.blobEngine.FindByBucket(bucket)
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

	if id == "" && bucket == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	blobInfo, err := bh.blobEngine.FindUnique(bucket, id)
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

	if bucket == "" || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	blobBytes, err := bh.blobEngine.Download(bucket, id)
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

	if bucket == "" && id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "verify the bucket or id sent"})
		return
	}

	_, err := bh.blobEngine.Delete(bucket, id)
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

	if _, err := bh.blobEngine.FindUnique(req.Bucket, req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "verify the data you sent and try again",
		})
		return
	}

	now := time.Now()
	signature := fmt.Sprintf("%s%s", shared.GenShortID(), now.Format("050204")) // format to SSDDMM

	bh.signaturesCache.Set(
		signature,
		domain.Signature{
			Bucket: req.Bucket,
			ID:     req.ID,
			TTL:    now.Add(time.Duration(req.TTL) * time.Minute),
			Once:   req.Once,
		},
	)

	c.JSON(http.StatusCreated, gin.H{
		"url": fmt.Sprintf("http://localhost:7713/b?sign=%s", signature),
	})
}
