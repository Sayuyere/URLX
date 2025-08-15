package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"urlx/shortener"
	"urlx/store"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Short string `json:"short"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var s store.Store = store.NewMemoryStore()
	var shortenerSvc shortener.Shortener = shortener.NewSimpleShortener()

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	r.POST("/shorten", func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
			c.Status(400)
			return
		}
		short := shortenerSvc.Shorten(req.URL)
		s.Set(short, req.URL)
		c.JSON(200, ShortenResponse{Short: short})
	})

	r.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")
		if long, ok := s.Get(short); ok {
			c.Redirect(302, long)
			return
		}
		c.Status(404)
	})

	r.DELETE("/delete/:short", func(c *gin.Context) {
		short := c.Param("short")
		s.Delete(short)
		c.Status(204)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	r.Run(":" + port)
}
