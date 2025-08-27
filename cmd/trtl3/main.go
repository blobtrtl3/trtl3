package main

import (
	"github.com/blobtrtl3/trtl3/api/handler"
	"github.com/blobtrtl3/trtl3/internal/db"
	"github.com/blobtrtl3/trtl3/internal/usecase/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	conn := db.NewDbConn()
	defer conn.Close()

	st := storage.NewBS(conn)

	bh := handler.NewBlob(st)

	r.POST("/blob", bh.Save)
	r.GET("/blob", bh.FindByID)
	r.GET("/blob", bh.FindByBucket)
	r.DELETE("/blob", bh.Delete)

	r.Run()
}
